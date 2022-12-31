package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// 理论上是初始奖励金额（但这里并未用算法对挖矿奖励进行处理，所以一直是这个值）
const SUBSIDY = 24

// 交易模型
type Transaction struct {
	TxHash []byte
	Vin    []TXinput
	Vout   []TXoutput
}

// 交易：序列化Transactions，用于打包交易得到交易hash
func (this *Transaction) Serialize() []byte {
	// 可供读写的缓存
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(this)
	if err != nil {
		fmt.Println("Fail to serialize transactions: ", err)
	}
	return res.Bytes()
}

// 验证是否是Coinbase交易
func (this *Transaction) IsCoinbase() bool {
	return len(this.Vin) == 1 && len(this.Vin[0].FromTxHash) == 0 && this.Vin[0].IndexOfVout == -1
}

// 生成TX hash
func (this *Transaction) SetTxHash() []byte {
	var resHash [32]byte

	// 将出去Hash其他信息序列化打包
	tx := *this
	tx.TxHash = []byte{}

	resHash = sha256.Sum256(tx.Serialize())

	return resHash[:]
}

// 对一笔交易签名
// 为了对一笔交易进行签名，需要获取交易输入所引用的输出，因为我们需要存储这些输出的交易
func (this *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if this.IsCoinbase() {
		return
	}
	// 在比特币里，所签名的并不是完整交易，而是一个去除部分Vin内容的交易副本
	txCopy := this.TrimmedCopy()
	// 遍历副本Vin，对每个输入分开签名
	for idxVin, vin := range txCopy.Vin {
		// 拿到输入对应的之前的UTXO
		prevTx := prevTXs[hex.EncodeToString(vin.FromTxHash)]
		// 双重检验: 每个输入中，Signature 被设置为 nil ，PubKey 被设置为所引用输出的 PubKeyHash
		txCopy.Vin[idxVin].Signature = nil
		txCopy.Vin[idxVin].PubKey = prevTx.Vout[vin.IndexOfVout].PubKeyHash
		// 对交易进行序列化，哈希后的结果就是我们要签名的数据
		txCopy.TxHash = txCopy.SetTxHash()
		txCopy.Vin[idxVin].PubKey = nil
		// 椭圆曲线签名
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic("Ecdsa.Sign Error!", err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		this.Vin[idxVin].Signature = signature
	}
}

// 验证一笔交易的签名
func (this *Transaction) Verify(prevTXs map[string]Transaction) bool {
	// 首先，需要同一笔交易的副本和同样算法曲线
	txCopy := this.TrimmedCopy()
	curve := elliptic.P256()

	// 检查每个输入的签名
	for idxVin, vin := range this.Vin {
		// 这部分与Sign相同，因为验证需要相同的数据
		prevTx := prevTXs[hex.EncodeToString(vin.FromTxHash)]
		txCopy.Vin[idxVin].Signature = nil
		txCopy.Vin[idxVin].PubKey = prevTx.Vout[vin.IndexOfVout].PubKeyHash
		txCopy.TxHash = txCopy.SetTxHash()
		txCopy.Vin[idxVin].PubKey = nil
		// 解包存储在Signature和PubKey中的值，在 crypto/ecdsa函数中使用。
		// 一个签名就是一对数字
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		// 一个公钥就是一对坐标
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		// 椭圆曲线算法回算公钥Hash, 再通过回算的公钥Hash签名验证
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}
	return true
}

// 获得修剪后交易副本
func (this *Transaction) TrimmedCopy() Transaction {
	var inputs []TXinput
	var outputs []TXoutput

	for _, vin := range this.Vin {
		// 裁剪Signature, PubKey
		inputs = append(inputs, TXinput{vin.FromTxHash, vin.IndexOfVout, nil, nil})
	}
	for _, vout := range this.Vout {
		outputs = append(outputs, TXoutput{vout.Value, vout.PubKeyHash})
	}
	txcopy := Transaction{
		// Hash 不变，部分裁剪input
		TxHash: this.TxHash,
		Vin:    inputs,
		Vout:   outputs,
	}
	return txcopy
}

// 交易 -格式化交易内容
func (this Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", this.TxHash))

	for i, input := range this.Vin {

		lines = append(lines, fmt.Sprintf("  Input %d:", i))
		lines = append(lines, fmt.Sprintf("    TXhash:      %x", input.FromTxHash))
		lines = append(lines, fmt.Sprintf("    FromOut:       %d", input.IndexOfVout))
		lines = append(lines, fmt.Sprintf("    Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("    PubKey:    %x", input.PubKey))
	}

	for i, output := range this.Vout {
		lines = append(lines, fmt.Sprintf("  Output %d:", i))
		lines = append(lines, fmt.Sprintf("    Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("    PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

// 创建Coinbase交易(Coinbase交易为挖矿奖励，没有输入，只有输出“凭空造币”，其data任意)
func NewCoinbaseTransaction(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXinput{
		FromTxHash:  []byte{},
		IndexOfVout: -1,
		Signature:   nil,
		PubKey:      []byte(data),
	}
	txout := NewTXOutput(SUBSIDY, to)

	tx := &Transaction{
		Vin:  []TXinput{txin},
		Vout: []TXoutput{*txout},
	}
	tx.TxHash = tx.SetTxHash()

	return tx
}

// 创建UTXO交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXinput
	var outputs []TXoutput

	// 创建并找到钱包账户的UTXO.outputs及其value
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	// DEBUG:fmt.Println(from)
	myWallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(myWallet.PublicKey)
	// DEBUG:fmt.Println("Sender-", fmt.Sprintf(" PubKeyHash:%x", pubKeyHash))

	balance, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if balance < amount {
		log.Panic("Error: No enough balance in your account")
	}

	// 创建inputs: 对应此次交易涉及到的有效output（因为此些output将被花费而失效）
	for txhash, outs := range validOutputs {
		// 拿到outputs hex hash []byte编码值
		txHash, err := hex.DecodeString(txhash)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			// 创建对应inputs
			input := TXinput{
				FromTxHash:  txHash,
				IndexOfVout: out,
				Signature:   nil,
				PubKey:      myWallet.PublicKey,
			}
			inputs = append(inputs, input)
		}
	}
	// 创建两个outputs：1.接受者地址锁定，实际转账 2. 发送者地址锁定，用于找零
	outputs = append(outputs, *NewTXOutput(amount, to))
	if amount < balance {
		outputs = append(outputs, *NewTXOutput(balance-amount, from))
	}

	tx := Transaction{
		Vin:  inputs,
		Vout: outputs,
	}
	tx.TxHash = tx.SetTxHash()

	// 最后，Sender用私钥对交易签名
	bc.SignTransaction(&tx, myWallet.PrivateKey)

	return &tx
}
