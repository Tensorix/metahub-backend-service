package onebot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type MessageResponse struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		MessageID int64   `json:"message_id"`
		Time      float64 `json:"time"`
	} `json:"data"`
	Message string `json:"message"`
}

func (bot *Onebot) SendToFriend(friendID int64, messages []Message) (int64, error) {
	var mr MessageResponse
	var friend Friend
	if !bot.Avaliable() {
		return 0, errors.New("bot is not avaliable")
	}

	action := ActionRequest{
		Action: "send_msg",
		Params: ActionParams{
			MessageType: "private",
			Message:     messages,
		},
	}
	DB.First(&friend, "id = ?", friendID)
	userid := friend.UID
	action.Params.UserID = userid

	data, err := json.Marshal(action)
	if err != nil {
		return 0, err
	}
	bot.mutex.Lock()
	bot.conn.WriteMessage(websocket.TextMessage, data)
	<-bot.msgSignal
	bot.mutex.Unlock()
	json.Unmarshal(bot.message, &mr)
	if mr.Retcode != 0 {
		return 0, errors.New("mr.Retcode != 0")
	}
	friendMessage := FriendMessage{
		MessageID:   mr.Data.MessageID,
		FriendID:    friendID,
		MessageTS:   time.Now().Unix(),
		SelfMessage: true,
		ReadMark:    true,
		Hide:        false,
		Revoke:      false,
	}
	DB.Create(&friendMessage)
	for _, message := range messages {
		friendSubMessage := FriendSubMessage{
			FriendMessageID: friendMessage.ID,
			IsText:          message.Type == "text",
			Message:         message.Data.Text,
		}
		DB.Create(&friendSubMessage)
	}
	return mr.Data.MessageID, nil
}
