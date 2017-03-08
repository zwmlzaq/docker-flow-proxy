package main

import (
	"github.com/ziutek/syslog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type handler struct {
	*syslog.BaseHandler
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(100, nil, false)}
	go h.mainLoop()
	return &h
}

func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		log.Printf("HAPRoxy: %s%s\n", m.Tag, m.Content)
	}
	h.End()
}

func startSyslog() {
	s := syslog.NewServer()
	s.AddHandler(newHandler())
	s.Listen("127.0.0.1:1514")

	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	s.Shutdown()
}

func main() {
	go startSyslog()
	NewArgs().Parse()
}
