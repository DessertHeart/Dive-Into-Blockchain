package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
)

// 挖矿难度（前24位 = Hex前6位 = 前3个字节）
const BITS = 24

// 针对每个区块的工作量证明
type PoW struct {
	block *Block
	targetBits *big.Int
}

// 打包准备验证的数据(block + BITS + nonce)
func (this *PoW) PackData(nonce int) []byte {
	data := bytes.Join(
		[][]byte {
			this.block.PrevBlockHash,
			this.block.HashTransactions(),
			[]byte(strconv.FormatInt(this.block.Timestamp, 10)),
			[]byte(strconv.FormatInt(int64(BITS), 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		},
		[]byte{},
	)
	return data
}

// Pow挖矿
func (this *PoW) Mine() (int, []byte) {
	// 大数比较与运算用此类型
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	// 防止nonce上溢
	maxNonce := math.MaxInt64

	// PoW核心逻辑
	fmt.Printf("Mining block with transactions: \"%s\"\n", this.block.Transactions)
	for {
		if nonce < maxNonce {
			data := this.PackData(nonce)
			hash = sha256.Sum256(data)
			hashInt.SetBytes(hash[:])
			// hashInt < targetBits（挖矿难度上界）时
			if hashInt.Cmp(this.targetBits) == -1 {
				fmt.Printf("%x", hash)
				break
			}else {nonce++}
		}else {
			log.Panic("nonce溢出，算力不足")
			break
		}
	}
	fmt.Printf("\n\n")
	// 重点，一定注意对于[]byte返回的范围[:]
	return nonce, hash[:]
}

// Pow工作量验证，用到挖矿得到的Nonce是否小于难度BITS
func (this *PoW) Validate() bool {
	var hashInt big.Int

	// 与挖矿区别就是，取最终通过的Nonce为入参
	data := this.PackData(this.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(this.targetBits) == -1
	return isValid
}

// 初始化PoW
func NewPoW(block *Block) *PoW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - BITS))
	pow := &PoW{
		block:      block,
		targetBits: target,
	}

	return pow
}