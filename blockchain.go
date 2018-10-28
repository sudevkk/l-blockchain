package main

import (
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

func (bc *Blockchain) AddBlock(data string) {
	block := NewBlock(data, bc.tip)
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
			block := NewGenesysBlock()
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
