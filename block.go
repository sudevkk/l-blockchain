package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Data          []byte
	Nounce        int64
	Hash          []byte
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)
	return result.Bytes()
}

func NewGenesysBlock() *Block {
	return NewBlock("Genesys Block", []byte{})
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{Timestamp: time.Now().UTC().Unix(), Data: []byte(data), PrevBlockHash: prevBlockHash}
	pow := MakeNewPOW(b)
	pow.Mine()
	return pow.b
}

func IntToHex(i int64) []byte {
	return []byte(strconv.FormatInt(int64(i), 10))
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
