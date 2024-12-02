package notifypage

import (
	"google.golang.org/grpc"

	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

var bots []*onebot.Onebot

type server struct {
	notify.UnimplementedNotifyServiceServer
}

func Register(s *grpc.Server) {
	notify.RegisterNotifyServiceServer(s, &server{})
	bots = authpage.Bots
}
