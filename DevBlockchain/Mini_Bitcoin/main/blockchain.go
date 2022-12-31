package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"github.com/boltdb/bolt"
)

// 创世区块coinbase data
const GENESISDATA = "The Times 17/Apr/2022 Dazso keep going further"

// 区块链 -切片是有序，存储有序的哈希；Map无序，存储hash -> block
type Blockchain struct {
	// 最后区块的hash
	TipHash []byte
	// boltdb数据库
	DB *bolt.DB
}

// 验证DB数据库是否存在
func dbExists() bool {
	if _, err := os.Stat("blockBucket"); os.IsNotExist(err) {
		return false
	}

	return true
}

// 打包新区块
func (this *Blockchain) MineBlock(txs []*Transaction) *Block {
	var lastHash []byte
	// 验证交易
	for _, tx := range txs {
		if this.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	// （只读方式）获取最后区块的hash（prevHash)
	err := this.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blockBucket"))
		lastHash = bucket.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(txs, lastHash)
	// 上链
	err = this.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blockBucket"))

		err = bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic("l -> hash failed in boltdb")
		}

		// 更新最后区块的Hash
		this.TipHash = newBlock.Hash
		return nil
	})

	return newBlock
}

// 初始化新链 -DB（存储核心）
func CreateBlockChain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tipHash []byte

	// address将接收挖出创世区块的奖励
	cbtx := NewCoinbaseTransaction(address, GENESISDATA)
	genesis := GenesisBlock(cbtx)

	db, err := bolt.Open("myBlockchain.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	// 创建表(Update:读写模式）
	err = db.Update(func(tx *bolt.Tx) error {
		// 初始化区块链
		// BoltDB数据库
		// block bucket键值对存储(byte array机制), 内容如下
		// 1） 32字节block.hash -> block 结构
		// 2） l -> 链中最后一个块的hash(tiphash)
		bucket, err := tx.CreateBucket([]byte("blockBucket"))
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tipHash = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	blockchain := &Blockchain{
		TipHash: tipHash,
		DB:      db,
	}
	return blockchain
}

// 显示 -创建已存在的链副本
func NewBlockchain() *Blockchain {
	if dbExists() {
		fmt.Println("No existing blockchain, plz create one")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open("myBlockchain.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blockBucket"))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// 交易 -使用私钥签名一包交易
func (this *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	// 循环调用Sign
	for _, vin := range tx.Vin {
		prevTx, err := this.FindTransaction(vin.FromTxHash)
		if err != nil {
			log.Panic("FindTransaction Error :", err)
		}
		prevTXs[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}
	tx.Sign(privateKey, prevTXs)
}

//交易 -验证一包交易
func (this *Blockchain) VerifyTransaction(tx *Transaction) bool {
	// 如果是Coinbase, 因为无Vin, 直接通过
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)
	// 循环调用Verify
	for _, vin := range tx.Vin {
		prevTx, err := this.FindTransaction(vin.FromTxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}
	return tx.Verify(prevTXs)
}
