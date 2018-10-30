package main

import (
	"fmt"
)

var subsidy int = 1

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

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

func (tx Transaction) SetId() {
	tx.ID = []byte("1234")
}

func (tx Transaction) prepareData() []byte {
	return []byte{}
}

func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetId()
	return &tx
}
