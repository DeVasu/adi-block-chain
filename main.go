package main

import (
	"fmt"
	"strconv"

	"github.com/DeVasu/adi-block-chain/blockchain"
)

func main() {

	chain := blockchain.InitBlockChain()
	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("second block after genesis")

	for idx, block := range chain.Blocks {
		fmt.Printf("Block Number #%d\n", idx)
		fmt.Printf("\tData: %s\n", block.Data)
		fmt.Printf("\tPrevious Hash: %x\n", block.PrevHash)

		pow := blockchain.NewProof(block)
		fmt.Printf("\tPoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Printf("\tHash: %x\n", block.Hash)
		fmt.Println()

	}

}
