package friendpage

import (
	"context"
	"log"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	friend "github.com/Tensorix/metahub-backend-service/gen/proto/v1/friend"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
)

func (s *server) Send(_ context.Context, in *friend.SendRequest) (*friend.SendResponse, error) {
	response := &friend.SendResponse{
		Result: auth.CheckResult_CHECK_RESULT_FAILED,
	}
	token := in.Token.Token
	username := authpage.GetUsername(token)
	if username == "" {
		return response, nil
	}
	response.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
	for _, bot := range onebot.Bots {
		if bot.Username != username || in.AccountTag != bot.AccountTag {
			continue
		}
		var messages []onebot.Message
		for _, message := range in.Messages {
			messageType := "text"
			if message.Type == friend.MessageType_MESSAGE_TYPE_IMAGE {
				messageType = "image"
			}
			messages = append(messages, onebot.Message{
				Type: messageType,
				Data: onebot.MessageData{
					Text: string(message.Content),
				},
			})
		}
		messageID, err := bot.SendToFriend(in.FriendId, messages)
		if err != nil {
			log.Println(err)
			return response, nil
		}
		response.Result = auth.CheckResult_CHECK_RESULT_SUCCESS
		response.MessageId = messageID
		break
	}
	return response, nil
}
