package net

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
)

type Buff struct {
	b       []byte
	buffLen int
}

func NewBuff(size int) *Buff {
	return &Buff{
		b: make([]byte, size),
	}
}

func (b *Buff) ReadAll(ctx context.Context, conn net.Conn) ([]byte, error) {
	b.Clear()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		buf, cont, err := b.Read(conn)
		if err != nil {
			return nil, err
		}
		if !cont {
			return buf, nil
		}
	}
}

func (b *Buff) getHeadAndPacketSize() (int, int, error) {
	if b.buffLen == 0 {
		return 0, 0, nil
	}
	bSize := int(b.b[0])
	headSize := 1
	if (bSize & 0x80) != 0 {
		if b.buffLen < 4 {
			return 0, 0, nil
		}
		headSize = 4
		bSize = int(binary.BigEndian.Uint32(b.b) & 0x7FFFFFFF)
		if bSize > readBuffLen {
			return 0, 0, errors.New("packet bigger buffer")
		}
	}
	return headSize, bSize, nil
}

func (b *Buff) Read(conn net.Conn) ([]byte, bool, error) {
	headSize, bSize, err := b.getHeadAndPacketSize()
	if err != nil {
		return nil, false, err
	}
	end := headSize + bSize
	if end == 0 {
		end = readBuffLen
	}

	rLen, err := conn.Read(b.b[b.buffLen:end])
	if err != nil {
		return nil, false, err
	}
	if rLen == 0 {
		return nil, true, nil
	}
	b.buffLen += rLen

	headSize, bSize, err = b.getHeadAndPacketSize()
	if err != nil {
		return nil, false, err
	}

	if b.buffLen < bSize+headSize {
		return nil, true, nil
	}

	if b.buffLen > bSize+headSize {
		return nil, false, errors.New("packet bigger size")
	}
	if bSize == 0 {
		b.Clear()
		return nil, false, nil
	}
	return b.b[headSize : bSize+headSize], false, nil
}

func (b *Buff) Clear() {
	b.buffLen = 0
}
