package ptr

import (
	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
)

func ToPointerPbUsers(pbUsers []pb.User) []*pb.User {
	res := make([]*pb.User, len(pbUsers))
	for i := range pbUsers {
		res[i] = &pbUsers[i]
	}
	return res
}
