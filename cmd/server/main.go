package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"pow/internal/logger"
	"pow/internal/net"
	"syscall"
)

func main() {
	logger.InitLog()

	log.Info("server start")
	defer log.Info("server end")

	addr := os.Getenv("addr")
	if len(addr) == 0 {
		addr = ":8080"
	}

	gr, appctx := errgroup.WithContext(context.Background())

	gr.Go(func() error {
		serv := net.Server{Addr: addr}
		return serv.Start(appctx)
	})

	gr.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		for {
			select {
			case <-appctx.Done():
				return nil
			case <-sigs:
				log.Info("caught stop signal, exiting")
				return context.Canceled
			}
		}
	})

	if err := gr.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
