package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"pow/internal/logger"
	"pow/internal/net"
)

func main() {
	logger.InitLog()

	log.Info("client start")
	defer log.Info("client end")

	addr := os.Getenv("addr")
	if len(addr) == 0 {
		addr = ":8080"
	}

	clnt := net.Client{Addr: addr}
	err := clnt.Connect(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	clnt.Close()
}
