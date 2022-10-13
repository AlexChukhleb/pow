package net

import (
	"crypto"
	"time"
)

const (
	keepAlive   = time.Second * 500
	readBuffLen = 1024
	maxLatency  = 10

	alg     = crypto.SHA1
	zerolen = 16
)

type State int8

const (
	StateHandshakeReq State = iota
	StateHandshakeResp
	StateHandshakeFinish
)
