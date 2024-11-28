package ws

import (
	"context"
	"time"
)

func (ws *WS) Shutdown() error {
	if !ws.Running {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ws.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
