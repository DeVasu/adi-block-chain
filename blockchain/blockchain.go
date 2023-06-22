package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

type BlockChain struct {
	// Blocks []*Block
	LastHash []byte
	Database *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleErr(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	HandleErr(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	HandleErr(err)
}

func InitBlockChain() *BlockChain {
	// return &BlockChain{[]*Block{Genesis()}}

	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Gensis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleErr(err)
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			HandleErr(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})

	HandleErr(err)

	return &BlockChain{lastHash, db}
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (iterator *BlockChainIterator) Next() *Block {

	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		HandleErr(err)
		var valCopy []byte
		valCopy, err = item.ValueCopy(nil)
		block = Deserialize(valCopy)
		return err
	})
	HandleErr(err)

	iterator.CurrentHash = block.PrevHash

	return block

}
