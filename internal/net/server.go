package net

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"pow/internal/quotes"
	"pow/pkg/pow"
)

type Server struct {
	Addr string
}

func (s *Server) Start(ctx context.Context) error {
	lconf := net.ListenConfig{KeepAlive: keepAlive}
	listener, err := lconf.Listen(ctx, "tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}
	log.WithField("address", listener.Addr().String()).Info("start listen")

	e := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			conn, err := listener.Accept()
			if err != nil {
				e <- err
				return
			}

			err = connSettings(conn)
			if err != nil {
				log.Error(err)
				_ = conn.Close()
				continue
			}

			go s.handler(ctx, conn)
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-e:
		return err
	}
}

func (s *Server) handler(ctx context.Context, conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	logLocal := log.WithField("addr", remoteAddr)

	defer func() {
		_ = conn.Close()
		logLocal.Info("close connection")
	}()

	logLocal.Info("incoming connection")

	buff := NewBuff(readBuffLen)

	var handshakeArr []byte

	state := StateHandshakeReq

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		b, cont, err := buff.Read(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logLocal.Error(err)
			return
		}
		if cont {
			continue
		}
		log.Info(len(b))

		switch state {
		case StateHandshakeReq:
			logLocal.Info("give handshake")
			handshakeArr = pow.Gen(alg, zerolen)
			if err := connWrite(conn, handshakeArr); err != nil {
				logLocal.Error(err)
				return
			}

			_ = conn.SetReadDeadline(time.Now().Add(time.Second * maxLatency))

			state = StateHandshakeResp

		case StateHandshakeResp:
			logLocal.Info("give handshake key ", hex.EncodeToString(b))
			ok, err := pow.Check(handshakeArr, b, maxLatency)
			if err != nil {
				logLocal.Error(err)
				return
			}
			if !ok {
				logLocal.Error("pow fail")
				return
			}

			if err := connWrite(conn, []byte(quotes.GetRandom())); err != nil {
				logLocal.Error(err)
				return
			}

			state = StateHandshakeFinish

		default:
			// TODO:
			return
		}

		buff.Clear()
	}
}
