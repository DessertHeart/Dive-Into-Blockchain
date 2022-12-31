package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 链上钱包存储结构
type Wallets struct {
	// 私钥和公钥都是随机的字节序列，无法打印显示，所以最终都是通过算法转为string（对应实际钱包中显示的地址）
	Wallets map[string]*Wallet
}

// 添加wallet到wallets
func (this *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	this.Wallets[address] = wallet

	return address
}

// 创建新钱包单元
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

// 存储 -存储钱包至地址
func (this *Wallets) LoadFromFile() error {
	walletFile := fmt.Sprintf("Wallet.dat")
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	this.Wallets = wallets.Wallets

	return nil
}

// 存储 -利用gob将Wallet存储到文件
func (this Wallets) SaveToFile() {
	var content bytes.Buffer
	walletFile := fmt.Sprintf("Wallet.dat")

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(this)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// 根据Address返回Wallet
func (this Wallets) GetWallet(address string) Wallet {
	return *this.Wallets[address]
}
