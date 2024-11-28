package ws

import (
	"net/http"

	"github.com/Tensorix/metahub-backend-service/onebot/wshandler"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WS struct {
	IP        string
	Port      int
	WSHandler *wshandler.WSHandler
	mux       *http.ServeMux
	server    *http.Server
	writer    http.ResponseWriter
	request   *http.Request
	hander    []func(conn *websocket.Conn, messageType int, message []byte) error
	Running   bool
	registed  bool
}
