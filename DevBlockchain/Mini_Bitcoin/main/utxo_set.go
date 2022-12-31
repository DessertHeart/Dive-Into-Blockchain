package main
// 数据（交易）分开存储：
// 实际交易被存储在区块链中 block bucket
// 未花费输出被存储在UTXO集中 chainstate bucket

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

// 通过公钥哈希，找到对应UTXO.output集 --余额查询
func (this UTXOSet) FindUTXO(pubKeyHash []byte) []TXoutput {
	var UTXOs []TXoutput
	db := this.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("ChainstateBucket"))
		// 光标：关联bucket（多用于遍历bucket)
		// 仅在打开Bucket后有效，关闭bucket后不可使用
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// 通过在UTXO集里查找对应Pubkeyhash的output ---转账
func (this UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := this.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ChainstateBucket"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

// 使用UTXO找到未花费输出，并存储到数据库 ---重新索引
func (this *UTXOSet) Reindex() {
	// chainstate 数据库
	// 32 字节的交易哈希 -> 该笔交易的未花费交易输出记录
	db := this.Blockchain.DB
	bucketName := []byte("ChainstateBucket")

	// 1.创建Bucket
	err := db.Update(func(tx *bolt.Tx) error {
		// 如果存在，先移除旧的（相当于每次交易都有变化，更新UTXO存储内容）
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			// 因为不存在而报错的情况不考虑
			log.Panic(err)
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// 2.找到所有未花费UTXO.output
	UTXOs := this.Blockchain.FindUTXO()
	// 3.将UTXO.output存入bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		for txHash, outs := range UTXOs {
			key, err := hex.DecodeString(txHash)
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// 更新UTXO集，使UTXO集时刻处于最新状态
func (this UTXOSet) UpdateUTXO(block *Block) {
	db := this.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("ChainstateBucket"))

		// 遍历区块的所有交易
		for _, tx := range block.Transactions {
			// 处理Vin，目的为处理旧的Vout(不包括coinbase)
			if tx.IsCoinbase() == false {
				// 遍历查询
				for _, vin := range tx.Vin {
					updateOuts := TXoutputs{}
					outsFromDB := bucket.Get(vin.FromTxHash)
					outs := DeserializeOutputs(outsFromDB)
					// 判断有无未花费output
					for idx, out := range outs.Outputs {
						if idx != vin.IndexOfVout {
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}
					// chainstate数据库操作
					if len(updateOuts.Outputs) == 0 {
						// 该笔交易中所有的output都被花费了
						err := bucket.Delete(vin.FromTxHash)
						if err != nil {
							log.Panic(err)
						}
					} else {
						err := bucket.Put(vin.FromTxHash, updateOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}

			}
			// 处理新的Vout, 直接添加
			newOutput := TXoutputs{}
			for _, out := range tx.Vout {
				newOutput.Outputs = append(newOutput.Outputs, out)
			}
			err := bucket.Put(tx.TxHash, newOutput.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
