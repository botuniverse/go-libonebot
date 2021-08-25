package main

import (
	"github.com/botuniverse/go-libonebot/comm"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	comm.StartHTTPTask("127.0.0.1", 8080)
}
