package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Transactions  []*Transaction
	Nounce        int64
	Hash          []byte
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)
	return result.Bytes()
}

func (b *Block) generateTransactionsHash() []byte {

	var idsJoined [][]byte
	for _, t := range b.Transactions {
		idsJoined = append(idsJoined, []byte(t.ID))
	}

	hash := sha256.Sum256(bytes.Join(idsJoined, []byte{}))
	return hash[:]
}

func NewGenesysBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func NewBlock(tx []*Transaction, prevBlockHash []byte) *Block {
	b := &Block{Timestamp: time.Now().UTC().Unix(), Transactions: tx, PrevBlockHash: prevBlockHash}
	pow := MakeNewPOW(b)
	pow.Mine()
	return pow.b
}

func IntToHex(i int64) []byte {
	return []byte(strconv.FormatInt(int64(i), 10))
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func Deserialize(blckData []byte) *Block {
	block := &Block{}
	reader := bytes.NewReader(blckData)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(block)
	if err != nil {
		log.Fatal(err)
	}
	return block
}
