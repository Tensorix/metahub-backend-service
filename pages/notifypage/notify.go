package notifypage

import (
	"google.golang.org/grpc"

	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
)

type server struct {
	notify.UnimplementedNotifyServiceServer
}

func Register(s *grpc.Server) {
	notify.RegisterNotifyServiceServer(s, &server{})
}
