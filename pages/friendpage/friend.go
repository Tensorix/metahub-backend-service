package friendpage

import (
	"google.golang.org/grpc"

	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
)

type server struct {
	friend.UnimplementedFriendServiceServer
}

func Register(s *grpc.Server) {
	friend.RegisterFriendServiceServer(s, &server{})
}
