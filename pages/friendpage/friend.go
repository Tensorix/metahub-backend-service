package friendpage

import (
	"google.golang.org/grpc"

	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

var bots []*onebot.Onebot

type server struct {
	friend.UnimplementedFriendServiceServer
}

func Register(s *grpc.Server) {
	friend.RegisterFriendServiceServer(s, &server{})
	bots = authpage.Bots
}
