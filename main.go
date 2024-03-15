package main

import (
	"git.countmax.ru/countmax/layoutconfig.api/infra"
	"github.com/sethvargo/go-signalcontext"
)

var (
	build   string
	githash string
	version = "1.1.21"
)

// @title LayoutConfig API
// @version 2.0/1.1.21
// @Description This is a general service for interacting with layout configuration.
// @Description Общий сервис для взаимодействия с конфигурацией схем размещения.

// @contact.name API Support
// @contact.url https://helpdesk.watcom.ru
// @contact.email 1020@watcom.ru

func main() {
	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel()
	serv := infra.NewServer(ctx, version, build, githash)
	serv.Run()
	<-ctx.Done()
	serv.Stop()
}
