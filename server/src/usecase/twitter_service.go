package usecase

import (
	"context"
	"errors"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/consts"
	"golang.org/x/xerrors"
)

type TwitterService interface {
	SearchUsers(ctx context.Context, query string) ([]domain.User, error)
	FollowUser(ctx context.Context, screenName string) error
	UnfollowUser(ctx context.Context, screenName string) error
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
func (s *twitterServiceImpl) SearchUsers(ctx context.Context, query string) ([]domain.User, error) {
	users := make([]domain.User, 0)

	page := 1
	for {
		searchedUsers, err := s.twitterClient.SearchUsers(ctx, query, page)
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

func (s *twitterServiceImpl) FollowUser(ctx context.Context, screenName string) error {
	if err := s.twitterClient.FollowUser(ctx, screenName); err != nil {
		return xerrors.Errorf("failed to twitterClient.FollowUser: %w", err)
	}
	return nil
}

func (s *twitterServiceImpl) UnfollowUser(ctx context.Context, screenName string) error {
	if err := s.twitterClient.UnfollowUser(ctx, screenName); err != nil {
		return xerrors.Errorf("failed to twitterClient.UnfollowUser: %w", err)
	}
	return nil
}
