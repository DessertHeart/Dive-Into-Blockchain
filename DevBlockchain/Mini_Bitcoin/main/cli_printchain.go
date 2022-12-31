package main

import (
	"fmt"
)

// CLI-打印区块命令
func (this *CLI) PrintChain() {
	bc := NewBlockchain()
	// 通过迭代器遍历显示
	bci := bc.Iterator()
	for {
		block := bci.Next()

		// 打印区块基本信息
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		// PoW验证
		pow := NewPoW(block)
		fmt.Printf("PoW Validate: %t\n", pow.Validate())
		// 打印交易信息
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		// 退出条件
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}