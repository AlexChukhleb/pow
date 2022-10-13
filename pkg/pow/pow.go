package pow

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"time"
	"unsafe"
)

func Check(arr []byte, key []byte, maxLatency int64) (bool, error) {
	algName, zeroLen, ts, err := unmarshal(arr)
	if err != nil {
		return false, err
	}

	dt := time.Now().Unix() - int64(ts)
	if dt < 0 || dt > maxLatency {
		return false, errors.New("timeout")
	}

	var h hash.Hash
	switch algName {
	case crypto.SHA1.String():
		if zeroLen > crypto.SHA1.Size()*8 {
			return false, errors.New("overbits")
		}
		h = sha1.New()

	case crypto.MD5.String():
		if zeroLen > crypto.MD5.Size()*8 {
			return false, errors.New("overbits")
		}
		h = md5.New()

	case crypto.SHA256.String():
		if zeroLen > crypto.SHA256.Size()*8 {
			return false, errors.New("overbits")
		}
		h = sha256.New()

	// TODO: other alg
	default:
		return false, errors.New("unknow alg")
	}

	h.Write(arr)
	h.Write(key)
	res := h.Sum(nil)

	return checkBufZeroBits(res, zeroLen), nil
}

func checkBufZeroBits(b []byte, zeroLen int) bool {
	i := 0
	for zeroLen > 8 {
		if b[i] != 0 {
			return false
		}
		zeroLen -= 8
		i++
	}

	if (b[i] >> (8 - zeroLen)) != 0 {
		return false
	}

	return true
}

func SearchKey(arr []byte, maxLatency int64) ([]byte, error) {
	var k uint64
	sz := int(unsafe.Sizeof(k))
	key := make([]byte, sz)

	tmr := time.NewTimer(time.Second * time.Duration(maxLatency))
	for {
		select {
		case <-tmr.C:
			//log.Info(hex.EncodeToString(key), "   ", len(key))
			return nil, errors.New("timeout")
		default:
		}

		binary.LittleEndian.PutUint64(key, k)

		kk := key[:]
		for len(kk) > 1 {
			lkk := len(kk) - 1
			if kk[lkk] != 0 {
				break
			}
			kk = kk[:lkk-1]
		}

		ok, err := Check(arr, kk, maxLatency)
		if err != nil {
			return nil, err
		}
		if ok {
			return kk, nil
		}
		k++
	}
}
