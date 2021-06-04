package interfaces

import (
	"context"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"golang.org/x/xerrors"
)

func FromPbUser(ctx context.Context, pbUser *pb.User) (domain.User, error) {
	if pbUser == nil {
		return domain.User{}, xerrors.Errorf("failed to convert pbUser to user")
	}
	return domain.User{
		ID:          pbUser.Id,
		ScreenName:  pbUser.ScreenName,
		Description: pbUser.Description,
	}, nil
}

func FromPbUsers(ctx context.Context, pbUsers []pb.User) ([]domain.User, error) {
	users := make([]domain.User, len(pbUsers))
	for i := range pbUsers {
		user, err := FromPbUser(ctx, &pbUsers[i])
		if err != nil {
			return nil, xerrors.Errorf("failed to FromPbUser: %w", err)
		}
		users[i] = user
	}
	return users, nil
}

func ToPbUser(ctx context.Context, user domain.User) pb.User {
	return pb.User{
		Id:          int64(user.ID),
		ScreenName:  user.ScreenName,
		Description: user.Description,
	}
}

func ToPbUsers(ctx context.Context, users []domain.User) []pb.User {
	pbUsers := make([]pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = ToPbUser(ctx, user)
	}
	return pbUsers
}
