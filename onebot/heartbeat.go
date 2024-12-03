package onebot

import (
	"encoding/json"
	"time"
)

type Heartbeat struct {
	Interval      int      `json:"interval"`
	MetaEventType string   `json:"meta_event_type"`
	Time          int      `json:"time"`
	SelfID        int      `json:"self_id"`
	PostType      string   `json:"post_type"`
}

func (bot *Onebot) heartbeat() error {
	var heartBeat Heartbeat
	err := json.Unmarshal(bot.message, &heartBeat)
	if err != nil {
		return err
	}

	bot.avaliableBefore = time.Now().Unix() + int64(heartBeat.Interval/1000)
	return nil
}
