package onebot

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

type FriendList struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    []struct {
		UserID   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Remark   string `json:"remark"`
	} `json:"data"`
}

type Friend struct {
	Id        uint
	AccountId uint
	Nickname  string
	UID       int64
	Remark    string
	Deleted   bool
}

func (bot *Onebot) GetFriendList() ([]Friend, error) {
	var fl FriendList
	var friends []Friend
	action := ActionRequest{
		Action: "get_friend_list",
	}
	if !bot.Avaliable() {
		return friends, errors.New("bot is not avaliable")
	}
	data, err := json.Marshal(action)
	if err != nil {
		return friends, err
	}
	bot.conn.WriteMessage(websocket.TextMessage, data)
	<-bot.msgSignal
	json.Unmarshal(bot.message, &fl)
	if err := DB.Where("account_id = ?", bot.AccountID).Find(&friends).Error; err != nil {
		log.Println(err)
	}
	for i := 0; i < len(friends); i++ {
		friends[i].Deleted = true
	}
	for i := 0; i < len(fl.Data); i++ {
		exist := false
		for j := 0; j < len(friends); j++ {
			if fl.Data[i].UserID == friends[j].UID {
				exist = true
				friends[j].Deleted = false
				break
			}
		}
		if !exist {
			friends = append(friends, Friend{
				AccountId: bot.AccountID,
				Nickname: fl.Data[i].Nickname,
				UID: fl.Data[i].UserID,
				Remark: fl.Data[i].Remark,
				Deleted: false,
			})
		}
	}
	DB.Save(friends)
	return friends, nil
}
