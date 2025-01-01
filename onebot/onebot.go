package onebot

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	DB   *gorm.DB
	Bots = make([]*Onebot, 0)
)

type Onebot struct {
	UserID          int32
	DockerID        uint
	AccountID       int32
	Username        string
	AccountTag      string
	IP              string
	Port            int
	mux             *http.ServeMux
	server          *http.Server
	writer          http.ResponseWriter
	request         *http.Request
	Running         bool
	realUID         int64
	registed        bool
	avaliableBefore int64
	conn            *websocket.Conn
	message         []byte
	msgSignal       chan struct{}
	FriendMessage   chan struct{}
	mutex           sync.Mutex
}

type MessageData struct {
	Text     string `json:"text"`
	File     string `json:"file"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Summary  string `json:"summary"`
	SubType  int    `json:"subType"`
}

type Message struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

type ActionParams struct {
	GroupID     int       `json:"group_id"`
	Message     []Message `json:"message"`
	AutoEscape  bool      `json:"auto_escape"`
	MessageType string    `json:"message_type"`
	MessageID   int       `json:"message_id"`
	ID          string    `json:"id"`
	UID         int64     `json:"user_id"`
	Times       int       `json:"times"`
	Duration    int       `json:"duration"`
	Enable      bool      `json:"enable"`
	Card        string    `json:"card"`
	GroupName   string    `json:"group_name"`
	IsDismiss   bool      `json:"is_dismiss"`
	Flag        string    `json:"flag"`
	Approve     bool      `json:"approve"`
	Remark      string    `json:"remark"`
	Nickname    string    `json:"nickname"`
}

type ActionRequest struct {
	Action string       `json:"action"`
	Params ActionParams `json:"params"`
}

func NewOnebot(username string, accountTag string, ip string, port int, userID int32, accountID int32) {
	bot := Onebot{
		UserID:        userID,
		AccountID:     accountID,
		Username:      username,
		AccountTag:    accountTag,
		IP:            ip,
		Port:          port,
		msgSignal:     make(chan struct{}),
		FriendMessage: make(chan struct{}),
	}
	bot.Register()
	bot.Run()
	Bots = append(Bots, &bot)
	log.Println("New bot:", bot.Username, bot.AccountTag, bot.IP, bot.Port)
}
