package interfaces

import (
	"context"
	"log"

	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/usecase"
	"golang.org/x/xerrors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type followHandlerImpl struct {
	pb.UnimplementedFollowServiceServer
	twitterService usecase.TwitterService
}

func NewFollowHandler(
	twitterService usecase.TwitterService,
) pb.FollowServiceServer {
	return &followHandlerImpl{
		twitterService: twitterService,
	}
}

func (h *followHandlerImpl) FollowUser(ctx context.Context, msg *pb.FollowUserMessage) (*emptypb.Empty, error) {
	log.Print("-----------FollowUser-----------------")
	if err := h.twitterService.FollowUser(ctx, msg.ScreenName); err != nil {
		return nil, xerrors.Errorf("failed to twitterService.FollowUser: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *followHandlerImpl) UnfollowUser(ctx context.Context, msg *pb.UnfollowUserMessage) (*emptypb.Empty, error) {
	log.Print("-----------UnfollowUser-----------------")
	return nil, nil
}
