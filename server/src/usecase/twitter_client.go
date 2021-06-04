package usecase

import (
	"context"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
)

type TwitterClient interface {
	SearchUsers(ctx context.Context, query string, page int) ([]domain.User, error)
	FollowUser(ctx context.Context, screenName string) error
}
