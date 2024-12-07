package notifypage

import (
	"context"
	"log"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"google.golang.org/grpc"
)

func (s *server) FriendMessage(in *auth.CheckRequest, stream grpc.ServerStreamingServer[notify.FriendMessageResponse]) error {
	token := in.Token
	username := authpage.GetUsername(token)
	if username == "" {
		response := notify.FriendMessageResponse{
			Result: auth.CheckResult_CHECK_RESULT_FAILED,
		}
		stream.Send(&response)
		return nil
	}
	ctx, close := context.WithCancel(context.Background())

	for _, bot := range onebot.Bots {
		if bot.Username != username {
			continue
		}
		go func() {
			ts := time.Now().Unix()
			for {
				select {
				case <-bot.FriendMessage:
					var friends []onebot.Friend
					if err := onebot.DB.Where("account_id = ?", bot.AccountID).Find(&friends).Error; err != nil {
						log.Println(err)
					}
					for _, f := range friends {
						var messages []onebot.FriendMessage
						if err := onebot.DB.Where("message_ts > ? AND friend_id = ?", ts, f.Id).Find(&messages).Error; err != nil {
							log.Println(err)
						}
						for _, msg := range messages {
							var subMessages []onebot.FriendSubMessage
							var notifyMessage []*notify.Message
							if err := onebot.DB.Where("friend_message_id = ?", msg.ID).Find(&subMessages).Error; err != nil {
								log.Println(err)
							}
							for _, subMsg := range subMessages {
								t := notify.MessageType_MESSAGE_TYPE_IMAGE
								if subMsg.IsText {
									t = notify.MessageType_MESSAGE_TYPE_TEXT
								}
								notifyMessage = append(notifyMessage, &notify.Message{
									Type: t,
									Text: subMsg.Message,
								})
							}
							err := stream.Send(&notify.FriendMessageResponse{
								SelfId:    bot.UID,
								MessageId: msg.MessageID,
								Timestamp: msg.MessageTS,
								Msg:       notifyMessage,
							})
							if err != nil {
								log.Println(err)
								close()
							}
						}
					}
					ts = time.Now().Unix()
				case <-ctx.Done():
					return
				case <- stream.Context().Done():
					close()
				}
			}
		}()
	}
	defer close()
	<-ctx.Done()
	return nil
}
