package notifypage

import (
	"time"

	protov1 "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	"google.golang.org/grpc"
)

var (
	interval = 5 * time.Second
)

func (s *server) Heartbeat(in *notify.HeartbeatRequest, stream grpc.ServerStreamingServer[notify.HeartbeatResponse]) error {
	response := notify.HeartbeatResponse{
		Result: protov1.CheckResult_CHECK_RESULT_FAILED,
		Details: []*notify.Detail{
			// {
			// 	Connected:  true,
			// 	AccountTag: "tag1",
			// },
		},
		Interval: int32(interval / time.Millisecond),
	}
	for {
		err := stream.Send(&response)
		if err != nil {
			return err
		}
		time.Sleep(interval)
	}
}
