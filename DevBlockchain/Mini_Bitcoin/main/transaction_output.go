package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go_prj22_Bitcoin/base58"
	"log"
)

// 交易输出 (每个UTXO是原子性的，交易只能全拿走再单独创建新的UTXO找零，输入同）
// -输出锁定/解锁 :注意，这里真实的比特币通过专门的脚本语言实现, 这里通过方法
type TXoutput struct {
	// satoshi的数量
	Value int
	// 公钥哈希，用于之后的比较(验证这笔output属于谁)
	PubKeyHash []byte
}

// 创建一个新TXoutput， 创建新交易时用
func NewTXOutput(value int, address string) *TXoutput {
	txo := &TXoutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// 锁定一个输出，用于给A地址转币时，解码得到A的公钥锁定以表示是其专属
func (this *TXoutput) Lock(address []byte) {
	pubKeyHash := base58.Base58Decode(address)
	// DEBUG:fmt.Println(fmt.Sprintf("Lock before pubKeyHash:%x", pubKeyHash))

	// version和checksum不需要
	pubKeyHash = pubKeyHash[len([]byte{VERSION}) : len(pubKeyHash) - 4]
	// DEBUG:fmt.Println(fmt.Sprintf("Lock pubKeyHash:%x", pubKeyHash))

	this.PubKeyHash = pubKeyHash
}

// UseKey辅助函数
// 通过hash检查，是否是对应的公钥提供的锁定
func (this *TXoutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, this.PubKeyHash) == 0
}

// TXoutput组
type TXoutputs struct {
	Outputs []TXoutput
}

// 序列化outputs，用于UTXO打包交易
func (this TXoutputs) Serialize() []byte {
	// 可供读写的缓存
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(this)
	if err != nil {
		fmt.Println("Fail to serialize transactions: ", err)
	}
	return res.Bytes()
}

// 反序列化解包[]TXoutputs
func DeserializeOutputs(data []byte) TXoutputs {
	var outputs TXoutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}