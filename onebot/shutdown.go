package onebot

import (
	"context"
	"time"
)

func (bot *Onebot) Shutdown() error {
	if !bot.Running {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := bot.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
