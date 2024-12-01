package onebot

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

type FriendList struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    []struct {
		UserID   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Remark   string `json:"remark"`
	} `json:"data"`
}

func (bot *Onebot) GetFriendList() (FriendList, error) {
	var fl FriendList
	action := ActionRequest{
		Action: "get_friend_list",
	}
	if !bot.Avaliable() {
		return fl, errors.New("bot is not avaliable")
	}
	data, err := json.Marshal(action)
	if err != nil {
		return fl, err
	}
	bot.conn.WriteMessage(websocket.TextMessage, data)
	<-bot.msgSignal
	json.Unmarshal(bot.message, &fl)
	return fl, nil
}
