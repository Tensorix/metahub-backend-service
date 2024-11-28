package ws

import (
	"log"
	"net/http"
	"strconv"
)

func (ws *WS) Register() {
	handler := ws.WSHandler
	ws.hander = append(
		ws.hander,
		handler.Heartbeat,
		handler.Lifecycle,
	)

	ws.mux = http.NewServeMux()
	ws.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.writer = w
		ws.request = r
		err := ws.multiHander()
		if err != nil {
			log.Println(err.Error())
		}
	})
	ws.server = &http.Server{
		Addr:    ws.IP + ":" + strconv.Itoa(ws.Port),
		Handler: ws.mux,
	}
	ws.registed = true
}

func (bot *WS) multiHander() error {
	conn, err := upgrader.Upgrade(bot.writer, bot.request, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		for i := 0; i < len(bot.hander); i++ {
			if err := bot.hander[i](conn, messageType, message); err != nil {
				return err
			}
		}
		log.Println(string(message))
	}
}
