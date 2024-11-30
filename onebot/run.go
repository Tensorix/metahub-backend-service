package onebot

func (ws *Onebot) Run() {
	if !ws.registed {
		return
	}
	go func() {
		ws.Running = true
		if err := ws.server.ListenAndServe(); err != nil {
			ws.Running = false
		}
	}()
}
