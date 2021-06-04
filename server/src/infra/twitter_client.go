package infra

import (
	"context"
	"log"

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

// pageは1が最小値、51が最大値
func (c *twitterClientImpl) Search(ctx context.Context, query string, page int) ([]domain.User, error) {
	if page < searchParamPageMinValue {
		return nil, consts.ErrTwitterSearchParamPagesTooSmall
	}
	if page > searchParamPageMaxValue {
		return nil, consts.ErrTwitterSearchParamPagesTooBig
	}

	twitterUsers, res, err := c.client.Users.Search(query, &twitter.UserSearchParams{
		Count: searchParamCount,
		Page:  page,
	})
	if res != nil {
		log.Printf("requested to [%s], response status code is [%d]", res.Request.URL, res.StatusCode)
	} else {
		log.Print("response of twitter users search api is nil")
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to call twitter user search api: %w", err)
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
