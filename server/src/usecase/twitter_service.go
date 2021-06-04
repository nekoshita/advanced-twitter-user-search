package usecase

import (
	"context"
	"errors"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/consts"
	"golang.org/x/xerrors"
)

type TwitterService interface {
	Search(ctx context.Context, query string) ([]domain.User, error)
}

type twitterServiceImpl struct {
	twitterClient TwitterClient
}

func NewTwitterService(
	twitterClient TwitterClient,
) TwitterService {
	return &twitterServiceImpl{
		twitterClient: twitterClient,
	}
}

// Search は最大で50回のAPIリクエストを行い、最大で1000件のユーザーを返す
func (s *twitterServiceImpl) Search(ctx context.Context, query string) ([]domain.User, error) {
	users := make([]domain.User, 0)

	page := 1
	for {
		searchedUsers, err := s.twitterClient.Search(ctx, query, page)
		if err != nil {
			switch {
			case errors.Is(err, consts.ErrTwitterSearchParamPagesTooBig):
				break
			default:
				return nil, xerrors.Errorf("failed to twitterClient.Search: %w", err)
			}
		}

		if len(searchedUsers) == 0 {
			break
		}

		users = append(users, searchedUsers...)

		page++
	}

	return users, nil
}
