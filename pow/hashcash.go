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

type (
	block struct {
		timestamp int64
		data      []byte
		hash      []byte
		nonce     int
	}

	hashCash struct {
		block      *block
		target     *big.Int
		targetBits uint
	}
)

func NewHashCash(timestamp int64, data string, targetBits uint) *hashCash {
	target := big.NewInt(1)
	return &hashCash{
		block:      &block{timestamp: timestamp, data: []byte(data), hash: []byte{}, nonce: 0},
		target:     target.Lsh(target, 256-targetBits),
		targetBits: targetBits,
	}
}

func (hc *hashCash) Calculate(ctx context.Context) (int, []byte, bool) {
	var (
		hashInt big.Int
		hash    [32]byte
		ok      bool
		nonce   = 0
	)

LOOP:
	for nonce < maxNonce {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}

		data := hc.prepareData(nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(hc.target) == -1 {
			ok = true
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:], ok
}

func (hc *hashCash) Verify(nonce int) bool {
	var hashInt big.Int

	data := hc.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(hc.target) == -1
}

func (hc *hashCash) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			hc.block.data,
			intToHex(hc.block.timestamp),
			intToHex(int64(hc.targetBits)),
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
