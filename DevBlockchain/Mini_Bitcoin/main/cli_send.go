package main

import (
	"fmt"
	"log"
)

// CLI-转账
func (this *CLI) Send(from, to string, value int) {
	if !ValidateAddress(from) {
		log.Panic("Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("Recipient address is not valid")
	}

	bc := NewBlockchain()
	utxoset := UTXOSet{bc}
	defer bc.DB.Close()

	// 这里把挖矿奖励给到转账操作的人，所以没笔交易签名都会加一笔coinbase交易
	tx := NewUTXOTransaction(from, to, value, bc)
	cbtx := NewCoinbaseTransaction(from, "")
	newblock := bc.MineBlock([]*Transaction{tx, cbtx})

	// 更新UTXO集
	utxoset.UpdateUTXO(newblock)

	fmt.Println("Transaction Success")
}
