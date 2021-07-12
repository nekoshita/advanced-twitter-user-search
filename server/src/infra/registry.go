package infra

import (
	"github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces"
	pb "github.com/nekoshita/advanced-twitter-user-search/server/src/interfaces/gen/advanced_twitter_user_search"
	"github.com/nekoshita/advanced-twitter-user-search/server/src/usecase"
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
	twitterClient := NewTwitterClient(
		twitterConsumerKey,
		twitterConsumerSecret,
		twitterUserAccessToken,
		twitterUserAccessSecret,
	)

	// services
	twitterService := usecase.NewTwitterService(twitterClient)

	// handler
	healthcheckHandler := interfaces.NewHealthServer()
	userHandler := interfaces.NewUserHandler(twitterService)
	followHandler := interfaces.NewFollowHandler(twitterService)

	// server
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthcheckHandler)
	pb.RegisterFollowServiceServer(server, followHandler)
	pb.RegisterUserServiceServer(server, userHandler)

	// add reflection for debug
	reflection.Register(server)

	return server
}
