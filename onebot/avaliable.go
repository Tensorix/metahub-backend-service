package onebot

import (
	"time"
)

func (bot *Onebot) Avaliable() bool {
	timestamp := time.Now().Unix()
	return timestamp < bot.avaliableBefore + 1
}
