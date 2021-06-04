package interfaces

import (
	"context"
	"log"

	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"google.golang.org/protobuf/types/known/emptypb"
)

type followHandlerImpl struct {
	pb.UnimplementedFollowServiceServer
}

func NewFollowHandler() pb.FollowServiceServer {
	return &followHandlerImpl{}
}

func (h *userHandlerImpl) FollowUser(context.Context, *pb.FollowUserMessage) (*emptypb.Empty, error) {
	log.Print("-----------FollowUser-----------------")
	return nil, nil
}

func (h *userHandlerImpl) UnfollowUser(context.Context, *pb.UnfollowUserMessage) (*emptypb.Empty, error) {
	log.Print("-----------UnfollowUser-----------------")
	return nil, nil
}
