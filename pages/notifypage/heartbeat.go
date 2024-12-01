package notifypage

import (
	"log"
	"time"

	protov1 "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"google.golang.org/grpc"
)

var interval = 5 * time.Second

func (s *server) Heartbeat(in *notify.HeartbeatRequest, stream grpc.ServerStreamingServer[notify.HeartbeatResponse]) error {
	bots := authpage.Bots
	token := in.Token
	username := authpage.GetUsername(token)
	for {
		var details []*notify.Detail
		for _, bot := range bots {
			if bot == nil {
				break
			}
			log.Println("Avaliable:",bot.Avaliable())
			log.Println("AccountTag:",bot.AccountTag)
			if bot.Username == username {
				detail := &notify.Detail{
					Connected:  bot.Avaliable(),
					AccountTag: bot.AccountTag,
				}
				details = append(details, detail)
			}
		}
		response := notify.HeartbeatResponse{
			Result:   protov1.CheckResult_CHECK_RESULT_SUCCESS,
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
