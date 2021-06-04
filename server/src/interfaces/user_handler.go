package interfaces

import (
	"context"
	"errors"
	"log"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/domain"
	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/consts"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/util/ptr"
	"golang.org/x/xerrors"
)

type userHandlerImpl struct {
	pb.UnimplementedUserServiceServer
	twitterClient TwitterClient
}

func NewUserHandler(
	twitterClient TwitterClient,
) pb.UserServiceServer {
	return &userHandlerImpl{
		twitterClient: twitterClient,
	}
}

func (h *userHandlerImpl) SearchUsers(ctx context.Context, msg *pb.SearchUsersMessage) (*pb.Users, error) {
	log.Print("-----------SearchUsers-----------------")
	log.Print(msg.Query)

	users := make([]domain.User, 0)

	page := 1
	for {
		searchedUsers, err := h.twitterClient.Search(ctx, msg.Query, page)
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

	pbUsers := ToPbUsers(ctx, users)

	return &pb.Users{
		Users: ptr.ToPointerPbUsers(pbUsers),
	}, nil
}
