package main

import (
	"bytes"
	_"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

// 区块
type Block struct {
	// 时间戳
	Timestamp int64
	// 区块哈希
	Hash []byte
	// 父哈希
	PrevBlockHash []byte
	// 交易
	Transactions []*Transaction
	// 随机数
	Nonce int
}

// 计算交易哈希（Block Header验证的一部分）
func (this *Block) HashTransactions() []byte {
	var txHashes [][]byte
	//var resHash [32]byte

	// 遍历所有交易的hash，最后打包为一个hash
	for _, tx := range this.Transactions {
		txHashes = append(txHashes, tx.TxHash)
	}
	// 返回的是[32]byte数组，需要通过[:]赋值给切片
	//resHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	// 通过Merkle树得到根节点hash
	merkleTree := NewMerkleTree(txHashes)

	return merkleTree.RootNode.HashData
}

// 存储：序列化this区块
func (this *Block) Serialize() []byte {
	// 可供读写的缓存
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(this)
	if err != nil {
		fmt.Println("Fail to serialize block: ", err)
	}
	return res.Bytes()
}

// 存储：反序列化this区块
func Deserialize(target []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(target))

	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Fail to deserialize block: ", err)
	}
	return &block
}

// 创建新区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Nonce:         0,
	}

	poW := NewPoW(block)
	nonce, hash := poW.Mine()

	// 重点，一定注意对于[]byte返回的范围[:]
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// 创世区块 -Coinbase交易
func GenesisBlock(coinbase *Transaction) *Block {
	gBlock := NewBlock([]*Transaction{coinbase}, []byte{})
	return gBlock
}
