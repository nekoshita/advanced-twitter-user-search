package interfaces

import (
	"context"
	"log"

	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/usecase"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/ptr"
	"golang.org/x/xerrors"
)

type userHandlerImpl struct {
	pb.UnimplementedUserServiceServer
	twitterService usecase.TwitterService
}

func NewUserHandler(
	twitterService usecase.TwitterService,
) pb.UserServiceServer {
	return &userHandlerImpl{
		twitterService: twitterService,
	}
}

func (h *userHandlerImpl) SearchUsers(ctx context.Context, msg *pb.SearchUsersMessage) (*pb.Users, error) {
	log.Print("-----------SearchUsers-----------------")
	log.Print(msg.Query)

	users, err := h.twitterService.SearchUsers(ctx, msg.Query)
	if err != nil {
		return nil, xerrors.Errorf("failed to twitterService.Search: %w", err)
	}

	pbUsers := ToPbUsers(ctx, users)

	return &pb.Users{
		Users: ptr.ToPointerPbUsers(pbUsers),
	}, nil
}
