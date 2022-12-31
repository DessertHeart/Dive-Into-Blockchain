package main

import (
	"fmt"
	"log"
)


//CLI-创建链
func (cli *CLI) createBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("Address is not valid")
	}

	bc := CreateBlockChain(address)
	defer bc.DB.Close()

	// 重新索引UTXO
	utxoset := UTXOSet{bc}
	utxoset.Reindex()

	fmt.Println("Create Blockchain success!")
}
