package notifypage

import (
	"context"
	"log"
	"os"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	notify "github.com/Tensorix/metahub-backend-service/gen/proto/v1/notify"
	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
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
						if err := onebot.DB.Where("message_ts > ? AND friend_id = ? AND self_message = 0", ts, f.Id).Find(&messages).Error; err != nil {
							log.Println(err)
						}
						for _, msg := range messages {
							var subMessages []onebot.FriendSubMessage
							var notifyMessage []*friend.Message
							if err := onebot.DB.Where("friend_message_id = ?", msg.ID).Find(&subMessages).Error; err != nil {
								log.Println(err)
							}
							for _, subMsg := range subMessages {
								t := friend.MessageType_MESSAGE_TYPE_TEXT
								content := []byte(subMsg.Message)
								if !subMsg.IsText {
									t = friend.MessageType_MESSAGE_TYPE_IMAGE
									log.Println("read", subMsg.Message)
									var err error
									content, err = os.ReadFile("cache/images/" + subMsg.Message)
									if err != nil {
										log.Println(err)
									}
								}
								notifyMessage = append(notifyMessage, &friend.Message{
									Type:    t,
									Content: content,
								})
							}
							err := stream.Send(&notify.FriendMessageResponse{
								Result:      auth.CheckResult_CHECK_RESULT_SUCCESS,
								FriendId:    msg.FriendID,
								SelfMessage: false,
								MessageId:   msg.MessageID,
								Timestamp:   msg.MessageTS,
								ReadMark:    false,
								Hide:        false,
								Revoke:      false,
								Messages:    notifyMessage,
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
				case <-stream.Context().Done():
					close()
				}
			}
		}()
	}
	defer close()
	<-ctx.Done()
	return nil
}
