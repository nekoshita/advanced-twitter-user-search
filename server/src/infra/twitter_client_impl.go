package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/consts"
	"github.com/nekoshita/go-twitter/twitter"
	"golang.org/x/xerrors"
)

type twitterClientImpl struct {
	client *twitter.Client
}

func NewTwitterClient(
	twitterConsumerKey string,
	twitterConsumerSecret string,
	twitterUserAccessToken string,
	twitterUserAccessSecret string,
) *twitterClientImpl {
	oauthConfig := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	token := oauth1.NewToken(twitterUserAccessToken, twitterUserAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	twitterClient := twitter.NewClient(httpClient)

	return &twitterClientImpl{
		client: twitterClient,
	}
}

// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/follow-search-get-users/api-reference/get-users-search
// 仕様で20が最大値
const (
	searchParamCount        = 20
	searchParamPageMinValue = 1
	searchParamPageMaxValue = 51
)

// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/follow-search-get-users/api-reference/get-users-search
// pageは1が最小値、51が最大値（つまり、最大で検索可能なユーザーは1000件まで）
func (c *twitterClientImpl) SearchUsers(ctx context.Context, query string, page int) ([]domain.User, error) {
	if page < searchParamPageMinValue {
		return nil, consts.ErrTwitterSearchParamPagesTooSmall
	}
	if page > searchParamPageMaxValue {
		return nil, consts.ErrTwitterSearchParamPagesTooBig
	}

	twitterUsers, resp, err := c.client.Users.Search(query, &twitter.UserSearchParams{
		Count: searchParamCount,
		Page:  page,
	})
	defer resp.Body.Close()
	c.logResponse(resp)
	if err != nil {
		return nil, xerrors.Errorf("failed to call twitter user search api: %w", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, consts.ErrTwitterApiUnauthorized
	}

	users := make([]domain.User, len(twitterUsers))
	for i, v := range twitterUsers {
		users[i] = domain.User{
			ID:          v.ID,
			ScreenName:  v.ScreenName,
			Description: v.Description,
		}
	}
	return users, nil
}

func (c *twitterClientImpl) FollowUser(ctx context.Context, screenName string) error {
	// すでにフォローしてる場合は200が返ってくる時と403が返ってくる場合がある
	user, resp, err := c.client.Friendships.Create(&twitter.FriendshipCreateParams{
		ScreenName: screenName,
	})
	defer resp.Body.Close()
	c.logResponse(resp)
	if err != nil {
		switch err.Error() {
		// 鍵アカウントをすでにフォローしてる場合、[twitter: 160 You've already requested to follow lemor303442.]というエラーが返ってくる
		// その場合はこのアプリケーション内部ではエラーとして扱わない
		// エラーの文字列比較はしたくないが、go-twitterの実装上仕方がない
		case fmt.Sprintf("twitter: 160 You've already requested to follow %s.", screenName):
			log.Print(err)
		default:
			return xerrors.Errorf("failed to call twitter create follow api: %w", err)
		}
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return consts.ErrTwitterApiUnauthorized
	}

	log.Printf("successfully followed or sent follow request to user [@%s]", user.ScreenName)

	return nil
}

func (c *twitterClientImpl) UnfollowUser(ctx context.Context, screenName string) error {
	// 鍵アカウントのフォロー申請の取り消しはできない
	user, resp, err := c.client.Friendships.Destroy(&twitter.FriendshipDestroyParams{
		ScreenName: screenName,
	})
	defer resp.Body.Close()
	c.logResponse(resp)
	if err != nil {
		return xerrors.Errorf("failed to call twitter destroy follow api: %w", err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return consts.ErrTwitterApiUnauthorized
	}

	log.Printf("successfully unfollowed user [@%s]. if private account, could not sent unfollow request", user.ScreenName)

	return nil
}

func (c *twitterClientImpl) logResponse(resp *http.Response) {
	// go-twitterは内部でgithub.com/dghubble/slingを利用してHTTPリクエストを送っている
	// slingはresponse.Bodyはすでにcloseしているので、ここでresponse.Bodyを出力することはできない
	log.Printf("requested to [%s:%s], response status code is [%d]", resp.Request.Method, resp.Request.URL, resp.StatusCode)
}
