package onebot

import (
	"github.com/Tensorix/metahub-backend-service/onebot/ws"
	"github.com/Tensorix/metahub-backend-service/onebot/wshandler"
)

type Onebot struct {
	Websocket *ws.WS
	WSHandler *wshandler.WSHandler
}

func NewOnebot(ip string, port int) Onebot {
	handler := &wshandler.WSHandler{}
	_ws := &ws.WS{
		IP:        ip,
		Port:      port,
		WSHandler: handler,
	}

	bot := Onebot{
		Websocket: _ws,
		WSHandler: handler,
	}
	bot.Websocket.Register()
	return bot
}

func (bot *Onebot) Running() bool {
	return bot.Websocket.Running
}

func (bot *Onebot) Run() {
	bot.Websocket.Run()
}

func (bot *Onebot) Shutdown() {
	bot.Websocket.Shutdown()
}

func (bot *Onebot) Avaliable() bool {
	return bot.WSHandler.Avaliable()
}
