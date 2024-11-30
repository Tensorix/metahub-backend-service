package onebot

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	timeout = 5 * time.Second
)

type Onebot struct {
	Username        string
	IP              string
	Port            int
	mux             *http.ServeMux
	server          *http.Server
	writer          http.ResponseWriter
	request         *http.Request
	Running         bool
	AccountId       int64
	registed        bool
	Connected       bool
	avaliableBefore int64
	conn            *websocket.Conn
	message         []byte
	msgsignal       chan struct{}
}

type ActionRequest struct {
	Action string `json:"action"`
	Params struct {
		GroupID     int    `json:"group_id"`
		Message     string `json:"message"`
		AutoEscape  bool   `json:"auto_escape"`
		MessageType string `json:"message_type"`
		MessageID   int    `json:"message_id"`
		ID          string `json:"id"`
		UserID      string `json:"user_id"`
		Times       int    `json:"times"`
		Duration    int    `json:"duration"`
		Enable      bool   `json:"enable"`
		Card        string `json:"card"`
		GroupName   string `json:"group_name"`
		IsDismiss   bool   `json:"is_dismiss"`
		Flag        string `json:"flag"`
		Approve     bool   `json:"approve"`
		Remark      string `json:"remark"`
		Nickname    string `json:"nickname"`
	} `json:"params"`
}

func NewOnebot(username string, ip string, port int) *Onebot {
	bot := Onebot{
		Username: username,
		IP:       ip,
		Port:     port,
		msgsignal: make(chan struct{}),
	}
	bot.Register()
	return &bot
}
