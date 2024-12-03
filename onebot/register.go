package onebot

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (bot *Onebot) Register() {
	bot.mux = http.NewServeMux()
	bot.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bot.writer = w
		bot.request = r
		err := bot.multiHander()
		if err != nil {
			log.Println(err)
		}
	})
	bot.server = &http.Server{
		Addr:    bot.IP + ":" + strconv.Itoa(bot.Port),
		Handler: bot.mux,
	}
	bot.registed = true
}

type MessageData struct {
	PostType      string `json:"post_type"`
	MetaEventType string `json:"meta_event_type"`
}

func (bot *Onebot) multiHander() error {
	conn, err := upgrader.Upgrade(bot.writer, bot.request, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	bot.conn = conn

	for {
		var data MessageData
		_, message, err := conn.ReadMessage()
		bot.message = message
		if err != nil {
			return err
		}

		log.Println(string(message))
		json.Unmarshal(message, &data)
		switch data.PostType {
		case "meta_event":
			switch data.MetaEventType {
			case "lifecycle":
				log.Println("exec Lifecycle")
				err = bot.lifecycle()
			case "heartbeat":
				log.Println("exec Heartbeat")
				err = bot.heartbeat()
			}
		default:
			bot.msgSignal <- struct{}{}
		}
		if err != nil {
			log.Println(err)
		}
	}
}
