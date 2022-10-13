package net

import (
	"context"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"net"
	"pow/pkg/pow"
	"time"
)

const (
	dialTimeout = time.Second
)

type Client struct {
	Addr string
	conn net.Conn
}

func (c *Client) Connect(ctx context.Context) error {
	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: keepAlive,
	}

	var err error

	c.conn, err = dialer.DialContext(ctx, "tcp", c.Addr)
	if err != nil {
		return err
	}

	err = connSettings(c.conn)
	if err != nil {
		_ = c.conn.Close()
		return err
	}

	ctxLocal, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	buff := NewBuff(readBuffLen)

	log.Info("handshake")
	err = connWrite(c.conn, []byte{34})
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("read handshake")
	arr, err := buff.ReadAll(ctxLocal, c.conn)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("search key")
	key, err := pow.SearchKey(arr, maxLatency)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("send key ", hex.EncodeToString(key))
	err = connWrite(c.conn, key)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("read verification")
	arr, err = buff.ReadAll(ctxLocal, c.conn)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("give verification: ", string(arr))

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
