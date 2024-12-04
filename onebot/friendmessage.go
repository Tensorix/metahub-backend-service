package onebot

import (
	"encoding/json"
	"log"
)

type BotFriendMessage struct {
	MessageID int64 `json:"message_id"`
	UserID    int64 `json:"user_id"`
	Message   []struct {
		Type string `json:"type"`
		Data struct {
			Text     string `json:"text"`
			File     string `json:"file"`
			Filename string `json:"filename"`
			URL      string `json:"url"`
			Summary  string `json:"summary"`
			SubType  int    `json:"subType"`
		} `json:"data"`
	} `json:"message"`
	Time uint64 `json:"time"`
}

type FriendMessage struct {
	ID          uint
	MessageID   int64
	FriendID    int64
	MessageTS   uint64
	SelfMessage bool
	ReadMark    bool
	Hide        bool
	Revoke      bool
}

type FriendSubMessage struct {
	ID              uint
	FriendMessageID uint
	IsText          bool
	Message         string
}

func (bot *Onebot) friendMessage() error {
	var botmsg BotFriendMessage
	err := json.Unmarshal(bot.message, &botmsg)
	if err != nil {
		log.Println(err)
		return nil
	}
	var friend Friend
	if err := DB.First(&friend, "uid = ? AND account_id = ?", botmsg.UserID, bot.AccountID).Error; err != nil {
		return err
	}
	friendMessage := FriendMessage{
		MessageID:   botmsg.MessageID,
		FriendID:    friend.Id,
		MessageTS:   botmsg.Time,
		SelfMessage: false,
		ReadMark:    false,
		Hide:        false,
		Revoke:      false,
	}
	DB.Create(&friendMessage)
	var friendSubMessages []FriendSubMessage
	for i := 0; i < len(botmsg.Message); i++ {
		istext := true
		msg := botmsg.Message[i].Data.Text
		if botmsg.Message[i].Type != "text" {
			istext = false
			msg = botmsg.Message[i].Data.URL
		}
		friendSubMessages = append(friendSubMessages, FriendSubMessage{
			FriendMessageID: friendMessage.ID,
			IsText:          istext,
			Message:         msg,
		})
	}
	DB.Create(&friendSubMessages)
	select {
	case <-bot.FriendMessage:
	default:
	}
	go func() {
		bot.FriendMessage <- struct{}{}
	}()
	return nil
}
