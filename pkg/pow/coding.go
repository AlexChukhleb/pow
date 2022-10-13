package pow

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"
	"unsafe"
)

// Gen
// return <alg>_<zerolen>_<ts>_<salt>
func Gen(alg crypto.Hash, zeroLen uint8) []byte {
	b := bytes.NewBuffer(make([]byte, 0, 127))
	b.WriteString(alg.String())
	b.Write([]byte{0})
	_ = binary.Write(b, binary.LittleEndian, zeroLen)
	b.Write([]byte{0})
	_ = binary.Write(b, binary.LittleEndian, uint32(time.Now().Unix()))
	b.Write([]byte{0})
	addSalt(b)

	return b.Bytes()
}

func addSalt(b *bytes.Buffer) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	sz := int(unsafe.Sizeof(r.Uint64()))
	for b.Len() < b.Cap()-sz {
		v := r.Uint64()
		_ = binary.Write(b, binary.LittleEndian, v)
	}
}

func unmarshal(arr []byte) (algName string, zeroLen int, ts uint32, err error) {
	if len(arr) < 1 {
		err = errors.New("short")
		return
	}

	b := bytes.NewBuffer(arr)

	algName, err = b.ReadString(0)
	if err != nil {
		return
	}
	if len(algName) == 0 {
		err = errors.New("can't get algName")
		return
	}
	algName = algName[:len(algName)-1]

	zeroLen_, err := b.ReadByte()
	if err != nil {
		return
	}
	zeroLen = int(zeroLen_)

	_, err = b.ReadByte()
	if err != nil {
		return
	}

	err = binary.Read(b, binary.LittleEndian, &ts)
	if err != nil {
		return
	}

	return
}
