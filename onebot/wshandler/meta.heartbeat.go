package wshandler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Heartbeat struct {
	Interval      int      `json:"interval"`
	Status        struct{} `json:"status"`
	MetaEventType string   `json:"meta_event_type"`
	Time          int      `json:"time"`
	SelfID        int      `json:"self_id"`
	PostType      string   `json:"post_type"`
}

func (handler *WSHandler) Heartbeat(conn *websocket.Conn, messageType int, message []byte) error {
	var heartBeat Heartbeat
	err := json.Unmarshal(message, &heartBeat)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if heartBeat.PostType != "meta_event" || heartBeat.MetaEventType != "heartbeat" {
		return nil
	}
	log.Println("Heartbeat")

	handler.avaliableBefore = time.Now().Unix() + int64(heartBeat.Interval) + int64(timeout)
	return nil
}
