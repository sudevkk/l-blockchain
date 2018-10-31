package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

var subsidy int = 50

type TxInput struct {
	TxID []byte
	// index of the transaction referenced
	Vout      int
	ScriptSig string
}

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

type MarkedTxOutput struct {
	Value        int
	Index        int
	ScriptPubKey string
}

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetID()
	return &tx
}

func NewUTXTransation(from string, to string, amount int, UTXs map[string][]MarkedTxOutput) *Transaction {

	var txIn []TxInput
	var txOut []TxOutput
	var totalInputAmount int

	for txId, markedUTXs := range UTXs {
		// Todo Check can Unlock
		for _, in := range markedUTXs {
			txIdHex, err := hex.DecodeString(txId)
			if err != nil {
				log.Panic(err)
			}
			txIn = append(txIn, TxInput{txIdHex, in.Index, from})
			totalInputAmount = totalInputAmount + in.Value
		}
	}

	txOut = append(txOut, TxOutput{amount, to})

	if amount < totalInputAmount {
		txOut = append(txOut, TxOutput{totalInputAmount - amount, from})
	}

	tx := Transaction{nil, txIn, txOut}
	tx.SetID()
	return &tx
}

func (in *TxInput) CanUnlockOutputWith(address string) bool {
	return in.ScriptSig == address
}

func (out *TxOutput) CanBeUnlockedWith(address string) bool {
	return out.ScriptPubKey == address
}
