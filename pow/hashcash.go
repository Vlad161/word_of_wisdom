package pow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

type hashCash struct {
}

func NewHashCash() *hashCash {
	return &hashCash{}
}

func (hc *hashCash) Calculate(ctx context.Context, payload []byte, timestamp int64, targetBits uint) (int, []byte, bool) {
	var (
		hashInt   big.Int
		hash      [32]byte
		ok        bool
		rawTarget = big.NewInt(1)
		target    = rawTarget.Lsh(rawTarget, 256-targetBits)
		nonce     = 0
	)

LOOP:
	for nonce < maxNonce {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}

		data := hc.prepareData(payload, timestamp, targetBits, nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			ok = true
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:], ok
}

func (hc *hashCash) Verify(payload []byte, timestamp int64, targetBits uint, nonce int) bool {
	var (
		hashInt   big.Int
		rawTarget = big.NewInt(1)
		target    = rawTarget.Lsh(rawTarget, 256-targetBits)
	)

	data := hc.prepareData(payload, timestamp, targetBits, nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1
}

func (hc *hashCash) prepareData(payload []byte, timestamp int64, targetBits uint, nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			payload,
			intToHex(timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	_ = binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}
