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
// max value for searchParamCount is 20
const (
	searchParamCount        = 20
	searchParamPageMinValue = 1
	searchParamPageMaxValue = 51
)

// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/follow-search-get-users/api-reference/get-users-search
// page must be bigger than 0, less than 52
// (you can search for 1,000 users at most for the given query)
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
	// https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/follow-search-get-users/api-reference/post-friendships-create
	// > If the user is already friends with the user a HTTP 403 may be returned,
	// > though for performance reasons this method may also return a HTTP 200 OK message even if the follow relationship already exists
	user, resp, err := c.client.Friendships.Create(&twitter.FriendshipCreateParams{
		ScreenName: screenName,
	})
	defer resp.Body.Close()
	c.logResponse(resp)
	if err != nil {
		switch err.Error() {
		// If you already sent follow request to a private account,
		// returned error will be "twitter: 160 You've already requested to follow lemor303442."
		// This application will not handle it as error.
		// (Hadnling errors by comparing string value is not preferred, but do it because of go-twitter implementation)
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
	// This will not destory follow request to private account
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
	// go-twitter uses github.com/dghubble/sling to send HTTP requests.
	// You cant log response.Body because sling closes response.Body inside its library
	log.Printf("requested to [%s:%s], response status code is [%d]", resp.Request.Method, resp.Request.URL, resp.StatusCode)
}
