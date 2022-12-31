package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	_ "go_prj22_Bitcoin/ripemd160"
	"go_prj22_Bitcoin/base58"
	"log"
)

// 版本号
const VERSION = byte(0x00)

/* --------------------------------钱包---------------------------------- */
/* -------------------------------Address-------------------------------- */

// 单个钱包结构
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

// 公钥 & 私钥 构造函数(注意: 内部函数)
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 生成P-256椭圆曲线
	curve := elliptic.P256()
	// crypto/rand.Reader:强随机数生成：约0-10^77
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 公钥是椭圆曲线上的点，即X坐标和Y坐标的拼接
	publicKey := append(privateKey.X.Bytes(), privateKey.Y.Bytes()...)

	return *privateKey, publicKey
}

// 验证地址是否有效
func ValidateAddress(address string) bool {
	// 用户输入的address
	pubKeyHash := base58.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash) - 4:]

	// 实际的地址
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// 通过公钥获得一个钱包地址（具体流程可见登链收藏-地址）
func (this *Wallet) GetAddress() []byte {
	// STEP1: RIPEMD160(SHA256(PubKey))
	pubHash := HashPubKey(this.PublicKey)
	// DEBUG:fmt.Println("该钱包公钥hash：", fmt.Sprintf("%x", pubHash))
	// STEP2: hash加上算法版本前缀
	versionPayLoad := append([]byte{VERSION}, pubHash...)
	// DEBUG:fmt.Println("该钱包versionPayLoad：", fmt.Sprintf("%x", versionPayLoad))
	// STEP3: 计算校验和
	checkSum := checksum(versionPayLoad)
	// STEP4: 校验和附加到payload
	fullPayLoad := append(versionPayLoad, checkSum...)
	// STEP5: Base58算法组合编码生成[]byte类型address
	address := base58.Base58Encode(fullPayLoad)

	return address
}

// STEP1
func HashPubKey(pubKey []byte) []byte {
	pubHashSHA256 := sha256.Sum256(pubKey)

	ripemd160Hasher := crypto.RIPEMD160.New()
	_, err := ripemd160Hasher.Write(pubHashSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	pubHashRIPEMD160 := ripemd160Hasher.Sum(nil)

	return pubHashRIPEMD160
}

// STEP3
func checksum(payload []byte) []byte {
	hashRes1 := sha256.Sum256(payload)
	hashRes2 := sha256.Sum256(hashRes1[:])

	// 校验和取结果哈希的前四个字节
	return hashRes2[:4]
}

// 钱包构造函数
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := &Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
	return wallet
}
