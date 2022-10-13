package net

import (
	"encoding/binary"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
)

func connSettings(conn net.Conn) error {
	tConn, ok := conn.(*net.TCPConn)
	if !ok {
		return errors.New("not tcp conn. close")
	}
	if err := tConn.SetKeepAlivePeriod(keepAlive); err != nil {
		return err
	}
	if err := tConn.SetKeepAlive(true); err != nil {
		return err
	}
	if err := tConn.SetNoDelay(true); err != nil {
		return err
	}
	return nil
}

func connWrite(conn net.Conn, b []byte) error {
	l1, err := conn.Write(headFromLen(len(b)))
	if err != nil {
		return err
	}

	l2, err := conn.Write(b)
	if err != nil {
		return err
	}

	log.Debug("written ", l1+l2)

	return nil
}

func headFromLen(len int) []byte {
	if len <= 127 {
		return []byte{byte(len)}
	}
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(len))
	head[0] |= 0x80
	return head
}
