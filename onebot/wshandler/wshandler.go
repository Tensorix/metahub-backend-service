package wshandler

import "time"

var (
	timeout = 1000
)

type WSHandler struct {
	AccountId       int64
	Connected       bool
	avaliableBefore int64
}

func (handler *WSHandler) Avaliable() bool {
	timestamp := time.Now().Unix()
	return timestamp < handler.avaliableBefore
}
