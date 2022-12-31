package main

import (
	"fmt"
	"go_prj22_Bitcoin/base58"
	"log"
)

// CLI -余额查询
func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain()
	utxoset := UTXOSet{bc}
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := base58.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := utxoset.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}