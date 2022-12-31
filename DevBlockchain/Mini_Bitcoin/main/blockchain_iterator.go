package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/boltdb/bolt"
	"log"
)
// 显示-BoltDB.key迭代器(一个一个读取，又不用全部加载db)
type BlockchainIterator struct {
	CurHash []byte
	DB *bolt.DB
}

// 显示-迭代器：返回链中下一块，从链尾到链头
func (this *BlockchainIterator) Next() *Block {
	var block *Block

	// 获取最后区块的hash（prevHash)
	err := this.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blockBucket"))
		encodedBlock := bucket.Get(this.CurHash)
		block = Deserialize(encodedBlock)

		return nil
	})
	if err != nil {log.Panic("DB.View failed")}
	// 向前滚动
	this.CurHash = block.PrevBlockHash
	return block
}

// 显示-创建迭代器，通过db将迭代器与Blockchain（存储了数据库boltDB链接的Blockchain实例）链接
func (this *Blockchain) Iterator() *BlockchainIterator {
	bcIterator := &BlockchainIterator{
		CurHash: this.TipHash,
		DB:      this.DB,
	}
	return bcIterator
}

// 交易 -通过TxHash查找交易
func (this *Blockchain) FindTransaction(TxHash []byte) (Transaction, error) {
	// 查找：迭代器辅助
	bci := this.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.TxHash, TxHash) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("Transaction is not found")
}

// ★★★交易-找到所有含有未花费输出的交易(未花费：输出并没有被任何交易的输入引用)
func (this *Blockchain) FindUnspentTx(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	// 记录已花费的output: hash string账户 --> IndexOfoutput
	spentTXOs := make(map[string][]int)
	// 通过迭代器遍历查询
	bci := this.Iterator()

	for{
		block := bci.Next()
		//遍历该区块内所有交易
		for _, tx := range block.Transactions {
			txHash := hex.EncodeToString(tx.TxHash)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 该笔UTXO.output已经加入过Vin（即被花费）, 无意义, 直接下一次查询
				if spentTXOs[txHash] != nil {
					for _, spentout := range spentTXOs[txHash] {
						if spentout == outIdx {
							continue Outputs
						}
					}
				}
				// ！注意，因为是倒序遍历，所有先加output再加Input没问题，因为最后一个块的最后一个交易的output肯定没被花费
				// 未被花费, 先检查公钥, 将交易加入List
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}

				// 花费后，（每一笔都如下检查）将该交易中Vin的加入已花费清单, 创世区块不存在Vin
				if tx.IsCoinbase() == false {
					for _, vin := range tx.Vin {
						if vin.UseKey(pubKeyHash) {
							vinHash := hex.EncodeToString(vin.FromTxHash)
							spentTXOs[vinHash] = append(spentTXOs[vinHash], vin.IndexOfVout)
						}
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

// 基本与Blockchain.FindUnspentTransactions一样原理,不过返回的为output
// 找到所有未花费交易output出并返回输出集（去掉了已花费部分）
func (this *Blockchain) FindUTXO() map[string]TXoutputs {
	UTXO := make(map[string]TXoutputs)
	spentTXOs := make(map[string][]int)
	bci := this.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txHash := hex.EncodeToString(tx.TxHash)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 是否花掉
				if spentTXOs[txHash] != nil {
					for _, spentOutIdx := range spentTXOs[txHash] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				// ！注意，因为是倒序遍历，所有先加output再加Input没问题，因为最后一个块的最后一个交易的output肯定没被花费
				outs := UTXO[txHash]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txHash] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.FromTxHash)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.IndexOfVout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// ★交易-找到pubKeyHash公钥哈希对应的所有UTXO.outputs(即输出)
// 返回的为交易哈希 -> 该交易下的输出output
func (this *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	// string账户 --> IndexOfoutput
	unspentOutputs := make(map[string][]int)
	// 找到所有UTXO
	unspentTXs := this.FindUnspentTx(pubKeyHash)
	// DEBUG:看一下有什么
	//for _, tx := range unspentTXs {
	//	tx.String()
	//}
	// 可花费value总和
	aggregate := 0

Find:
	// 遍历UTXO
	for _, tx := range unspentTXs {
		txHash := hex.EncodeToString(tx.TxHash)

		for outIdx, out := range tx.Vout {
			// 判断是否为from的UTXO.output,
			if out.IsLockedWithKey(pubKeyHash) && aggregate < amount {
				aggregate += out.Value
				// DEBUG:fmt.Println("aggregate的值：", aggregate)
				unspentOutputs[txHash] = append(unspentOutputs[txHash], outIdx)

				// 找到满足此次转账金额即可停止
				if aggregate >= amount {
					break Find
				}
			}
		}
	}
	return aggregate, unspentOutputs
}