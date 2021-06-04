# twitter-advanced-user-search server

## installation

install go
```
$ brew install anyenv
$ brew install grpc
$ anyenv install goenv
$ goenv install 1.16.4
```

install cli
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
$ brew install grpcurl
```

install modules
```
$ go mod vendor
```

## Generate GRPC Client
```
$ bin/gen_grpc_client
```

## Debugging
```
$ grpcurl -plaintext localhost:50001 list
$ grpcurl -plaintext -d '{"query":"Java"}' localhost:50001 advanced_twitter_user_search.UserService/SearchUsers
```
