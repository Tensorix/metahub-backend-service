package wshandler

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Lifecycle struct {
	SubType       string `json:"sub_type"`
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfID        int64  `json:"self_id"`
	PostType      string `json:"post_type"`
}

func (handler *WSHandler) Lifecycle(conn *websocket.Conn, messageType int, message []byte) error {
	var lifecycle Lifecycle
	err := json.Unmarshal(message, &lifecycle)
	if err != nil {
		log.Println(err)
		return nil
	}
	if lifecycle.PostType != "meta_event" || lifecycle.MetaEventType != "lifecycle" {
		return nil
	}
	log.Println("Lifecycle")
	
	handler.AccountId = lifecycle.SelfID
	handler.Connected = true
	return nil
}
