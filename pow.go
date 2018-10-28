package main

import (
	"bytes"
	"crypto/sha256"
	//"fmt"
	"math"
	"math/big"
)

const targetBits = 24

type POW struct {
	b      *Block
	target *big.Int
}

func (pow *POW) generateHash(nounce int64) []byte {
	joined := [][]byte{IntToHex(pow.b.Timestamp), pow.b.PrevBlockHash, pow.b.Data, IntToHex(nounce)}
	hash := sha256.Sum256(bytes.Join(joined, []byte{}))
	return hash[:]
}

func MakeNewPOW(b *Block) *POW {
	var maxTarget *big.Int = big.NewInt(1)
	maxTarget.Lsh(maxTarget, 256-16)
	//fmt.Printf("Target is %x", maxTarget)
	return &POW{b, maxTarget}
}

func (pow *POW) Mine() []byte {

	var testHashNum *big.Int = big.NewInt(0)
	var nounce int64 = 1

	for nounce < math.MaxInt64 {
		testHash := pow.generateHash(nounce)
		testHashNum.SetBytes(testHash)

		if testHashNum.Cmp(pow.target) == -1 {
			pow.b.Hash = testHash
			pow.b.Nounce = nounce
			break
		}

		nounce = nounce + 1
	}

	return pow.b.Hash
}