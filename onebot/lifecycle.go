package onebot

import (
	"encoding/json"
	"log"
)

type Lifecycle struct {
	SubType       string `json:"sub_type"`
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfID        int64  `json:"self_id"`
	PostType      string `json:"post_type"`
}

func (bot *Onebot) Lifecycle() error {
	var lifecycle Lifecycle
	err := json.Unmarshal(bot.message, &lifecycle)
	if err != nil {
		log.Println(err)
		return nil
	}
	// log.Println("Lifecycle")

	bot.AccountId = lifecycle.SelfID
	bot.Connected = true
	return nil
}
