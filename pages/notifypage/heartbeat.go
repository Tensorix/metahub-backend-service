package notifypage

import (
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"google.golang.org/grpc"
)

var interval = 5 * time.Second

func (s *server) Heartbeat(in *auth.CheckRequest, stream grpc.ServerStreamingServer[notify.HeartbeatResponse]) error {
	token := in.Token
	username := authpage.GetUsername(token)
	if username == "" {
		response := notify.HeartbeatResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		}
		stream.Send(&response)
		return nil
	}
	for {
		var details []*notify.Detail
		for _, bot := range onebot.Bots {
			if bot.Username != username {
				continue
			}
			
			detail := &notify.Detail{
				Connected:  bot.Avaliable(),
				AccountTag: bot.AccountTag,
			}
			details = append(details, detail)
		}
		response := notify.HeartbeatResponse{
			Result:   auth.CheckResult_CHECK_RESULT_SUCCESS,
			Details:  details,
			Interval: int32(interval / time.Millisecond),
		}

		err := stream.Send(&response)
		if err != nil {
			return err
		}
		time.Sleep(interval)
	}
}
