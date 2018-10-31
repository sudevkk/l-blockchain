package main

import (
	"encoding/hex"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	db  *bolt.DB
	tip []byte
}

type Reader struct {
	bc      *Blockchain
	current []byte
}

func (r *Reader) Next() (*Block, error) {
	var blockData []byte
	err := r.bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockData = b.Get(r.current)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if blockData == nil {
		return nil, errors.New("Block Was not Found")
	}
	block := Deserialize(blockData)
	r.current = block.PrevBlockHash
	return block, nil
}

func (bc *Blockchain) NewReader() *Reader {
	return &Reader{bc, bc.tip}
}

func (bc *Blockchain) AddBlock(tx []*Transaction) {
	block := NewBlock(tx, bc.tip)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		err := b.Put(block.Hash, block.Serialize())
		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	bc.tip = block.Hash
}

func MakeNewBlockchain(fname string) *Blockchain {
	var lastHash []byte
	db, err := bolt.Open(fname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("blocks"))
		if err != nil {
			log.Fatal(err)
		}
		lastHash = b.Get([]byte("l"))
		if lastHash == nil {
			block := NewGenesysBlock(fname)
			lastHash = block.Hash
			err = b.Put(block.Hash, block.Serialize())
			err = b.Put([]byte("l"), block.Hash)
		} else {
			lastHash = b.Get([]byte("l"))
		}
		return nil
	})

	bc := &Blockchain{db, lastHash}
	return bc
}

func (bc *Blockchain) getBalance(address string) int {
	balance, _ := bc.FindUnSpendTransactions(address, -1)
	log.Println("****")
	/*
		for _, tx := range unspentTXOs {
			fmt.Println("----")
			for _, out := range tx {
				balance = balance + out.Value
			}
		} */

	return balance
}

func (bc *Blockchain) FindUnSpendTransactions(address string, stopAtAmount int) (int, map[string][]MarkedTxOutput) {

	var unspentTXOs = make(map[string][]MarkedTxOutput)
	var spenTxs = make(map[string][]int)
	var unusedAmountSoFar int
	iterator := bc.NewReader()

blockIterator:
	for {
		log.Println("Iteratin block")
		b, _ := iterator.Next()

		for _, tx := range b.Transactions {

			txidstring := hex.EncodeToString(tx.ID)
			log.Println("Iteratin Trans - ", txidstring)
			// Sum Unspent
		checkIfNextOutputUsed:
			for index, out := range tx.Vout {
				// If the We have requested amount, break and return
				if stopAtAmount != -1 && unusedAmountSoFar >= stopAtAmount {
					break blockIterator
				}

				log.Println("Iteratin Outputs")
				log.Println("Checking Lock", address, out)
				if out.CanBeUnlockedWith(address) {
					log.Println("Yes - can be unlocked  ", address)
					usages, isUsed := spenTxs[txidstring]

					if isUsed {
						log.Println("Trans has Used Inputs")
						for _, usedIndex := range usages {
							if usedIndex == index {
								log.Println("Used Inputs - IGNORING")
								continue checkIfNextOutputUsed
							}
						}
					} else {
						log.Println("Trans has NO Used Inputs")
					}

					log.Println("Yes - can be used  ")
					markedUTX := MarkedTxOutput{out.Value, index, out.ScriptPubKey}
					unspentTXOs[txidstring] = append(unspentTXOs[txidstring], markedUTX)
					unusedAmountSoFar = unusedAmountSoFar + out.Value
				}
			}

			// Map Used Outputs Used in Inputs
			if true { // not genesys block
				log.Println("Iteratin Inputs")
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {

						inputTXID := hex.EncodeToString(in.TxID)
						log.Println("Input Unlocked ", inputTXID, " - ", in.Vout)
						spenTxs[inputTXID] = append(spenTxs[inputTXID], in.Vout)
					}
				}
			}
		}

		if len(b.PrevBlockHash) == 0 {
			break
		}

	}
	log.Println(" returning from unspent  ", unusedAmountSoFar)
	return unusedAmountSoFar, unspentTXOs
}

func (bc Blockchain) transfer(from string, to string, amount int) (*Block, error) {
	balance, unspentTXOs := bc.FindUnSpendTransactions(from, amount)

	if balance < amount {
		return &Block{}, errors.New("Insufficient funds")
	}

	tx := NewUTXTransation(from, to, amount, unspentTXOs)

	bc.AddBlock([]*Transaction{tx})
	return &Block{}, nil
}
