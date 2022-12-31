package main

import (
	"bytes"
)

// 交易输入 -输入签名
type TXinput struct {
	// 输出交易哈希，一个输入（必须）引用之前交易的输出
	FromTxHash []byte
	// 对应输出的Index, 辅助验证
	IndexOfVout int
	// 交易签名
	Signature []byte
	// 原生公钥， 用于验证： 对应所引用output的PubKeyHash
	PubKey []byte
}

// 检查是否是用该input的公钥，来解锁该输出，对比input原生公钥Hash和outputs里存储的Hash
func (this *TXinput) UseKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(this.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}