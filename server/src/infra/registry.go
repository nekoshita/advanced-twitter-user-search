package infra

import (
	"github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces"
	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func NewServer(
	twitterConsumerKey string,
	twitterConsumerSecret string,
	twitterUserAccessToken string,
	twitterUserAccessSecret string,
) *grpc.Server {
	// client
	twitterClient := NewTiwtterClient(
		twitterConsumerKey,
		twitterConsumerSecret,
		twitterUserAccessToken,
		twitterUserAccessSecret,
	)

	// handler
	healthcheckHandler := interfaces.NewHealthServer()
	userHandler := interfaces.NewUserHandler(twitterClient)
	followHandler := interfaces.NewFollowHandler()

	// server
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthcheckHandler)
	pb.RegisterUserServiceServer(server, userHandler)
	pb.RegisterFollowServiceServer(server, followHandler)

	// add reflection for debug
	reflection.Register(server)

	return server
}
