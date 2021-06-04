package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nekoshita/advanced-twitter-user-search/server/src/infra"
)

type config struct {
	Port                    string
	TwitterConsumerKey      string
	TwitterConsumerSecret   string
	TwitterUserAccessToken  string
	TwitterUserAccessSecret string
}

const followTargetTwitterUserScreenName = "nekoshita_yuki"

func main() {
	// read credentials from environment variables if available
	conf := &config{
		Port:                    os.Getenv("PORT"),
		TwitterConsumerKey:      os.Getenv("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret:   os.Getenv("TWITTER_CONSUMER_SECRET"),
		TwitterUserAccessToken:  os.Getenv("TWITTER_USER_ACCESS_TOKEN"),
		TwitterUserAccessSecret: os.Getenv("TWITTER_USER_ACCESS_SECRET"),
	}
	// allow consumer credential flags to override confing fields
	port := flag.String("port", "", "Port")
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flag.String("access-token", "", "Twitter User Access Token")
	accessSecret := flag.String("access-secret", "", "Twitter User Access Secret")
	flag.Parse()
	if *port != "" {
		conf.Port = *port
	}
	if *consumerKey != "" {
		conf.TwitterConsumerKey = *consumerKey
	}
	if *consumerSecret != "" {
		conf.TwitterConsumerSecret = *consumerSecret
	}
	if *accessToken != "" {
		conf.TwitterUserAccessToken = *accessToken
	}
	if *accessSecret != "" {
		conf.TwitterUserAccessSecret = *accessSecret
	}
	if conf.TwitterConsumerKey == "" {
		log.Fatal("Missing Twitter Consumer Key")
	}
	if conf.TwitterConsumerSecret == "" {
		log.Fatal("Missing Twitter Consumer Secret")
	}
	if conf.TwitterUserAccessToken == "" {
		log.Fatal("Missing Twitter User Access Token")
	}
	if conf.TwitterUserAccessSecret == "" {
		log.Fatal("Missing Twitter User Access Secret")
	}

	server := infra.NewServer(
		conf.TwitterConsumerKey,
		conf.TwitterConsumerSecret,
		conf.TwitterUserAccessToken,
		conf.TwitterUserAccessSecret,
	)

	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("staring server at port :%s", conf.Port)
	log.Panic(server.Serve(listenPort))
}
