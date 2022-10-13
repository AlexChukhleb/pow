package net

import (
	"net"
	"testing"
)

func TestBuffEmpty(t *testing.T) {
	buff := NewBuff(readBuffLen)

	r, w := net.Pipe()

	arr := []byte{}
	go connWrite(w, arr)

	b, cont, err := buff.Read(r)
	if cont {
		b, cont, err = buff.Read(r)
	}
	if err != nil {
		t.Fatal(err)
	}
	if cont {
		t.Fatal("continue")
	}
	if len(b) != len(arr) {
		t.Fatal("in != out")
	}
}

func TestBuffSmall(t *testing.T) {
	buff := NewBuff(readBuffLen)

	r, w := net.Pipe()

	arr := []byte{0}
	go connWrite(w, arr)

	b, cont, err := buff.Read(r)
	if cont {
		b, cont, err = buff.Read(r)
	}
	if err != nil {
		t.Fatal(err)
	}
	if cont {
		t.Fatal("continue")
	}
	if len(b) != len(arr) {
		t.Fatal("in != out")
	}
}

func TestBuffBig(t *testing.T) {
	buff := NewBuff(readBuffLen)

	r, w := net.Pipe()

	arr := make([]byte, readBuffLen-4)
	go connWrite(w, arr)

	b, cont, err := buff.Read(r)
	if cont {
		b, cont, err = buff.Read(r)
	}
	if err != nil {
		t.Fatal(err)
	}
	if cont {
		t.Fatal("continue")
	}
	if len(b) != len(arr) {
		t.Fatal("in != out")
	}
}
