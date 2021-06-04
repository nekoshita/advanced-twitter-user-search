package usecase

import (
	"context"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
)

type TwitterClient interface {
	Search(ctx context.Context, query string, page int) ([]domain.User, error)
}
