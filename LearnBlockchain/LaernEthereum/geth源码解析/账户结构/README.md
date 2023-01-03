### 账户结构

> 以太坊基本数据单元(Metadata)：`Account`

不同于Bitcoin 基于 UTXO 的 Blockchain/Ledger 系统，Ethereum是以 **Account/State 模型为核心的基于交易驱动的状态机(Transaction-based State Machine)**，我们常说的世界状态(State)，其根本是反应了某一账户(Account)在某一时刻下的属性。

其中，State 对应的基本数据结构，称为 **stateObject**（内部维护，不可导出）。当 stateObject 的值由Transaction 的执行而触发数据更新/删除/创建而发生了变化时，我们称为**状态转移**。

> 即stateObject 的状态从当前的 State 转移到另一个 State

```go
// core/state/state_object.go
type stateObject struct {
  address  common.Address
  addrHash common.Hash // hash(address)
  
  // account state
  data     types.StateAccount
  // 指向stateDB: 真正存储数据的地方
  // 方便调用 StateDB 相关的API对Account所对应的stateObject进行CRUD操作
  db       *StateDB

  // DB error.
  // State objects are used by the consensus core and VM which are
  // unable to deal with database-level errors. Any error that occurs
  // during a database read is memoized here and will eventually be returned
  // by StateDB.Commit.
  dbErr error

  // 内存缓存相关逻辑
  trie Trie // storage trie
  code Code // contract bytecode, 缓存代码当代码被从DB storage中加载出来 

  // 在执行 Transaction 的时候缓存合约修改的持久化数据
  // EOA账户为空
  originStorage  Storage 
  pendingStorage Storage 
  dirtyStorage   Storage 
  fakeStorage    Storage 

  dirtyCode bool // true: 当code被更新
  suicided  bool
  deleted   bool
}
```

#### 1）EOA账户

的创建：包含**本地创建**和**链上注册**(stateDB进行链上账户管理)，入口函数`NewAccount`

> `passphrase`入参仅用于加密本地保存私钥的**Keystore**文件（使用对称加密算法来加密私钥生成），与生成账户的私钥、地址的生成无关。
>
> **私钥泄露、助记词泄露、Keystore+密码泄露**都会导致账户控制权丢失

```go
// accounts/keystore/keystore.go
// 该API是geth暴露出来，用于方便用户本地创建于管理账户的
func (ks *KeyStore) NewAccount(passphrase string) (accounts.Account, error) {
  // 生成account的函数(ECDSA)
  _, account, err := storeNewKey(ks.storage, crand.Reader, passphrase)
  if err != nil {
    return accounts.Account{}, err
  }

  // 缓存并等待系统落盘
  ks.cache.add(account)
  ks.refreshWallets()
  return account, nil
}
```

- **[账号创建](https://zhuanlan.zhihu.com/p/53827188)：**

  > 这里实际上只进行了ecdsa计算和keystore存储，实际账户在以太坊世界状态中存储，只需要等待有相关联的Transaction发生，若不存在就会自动通过`newObject()  //core/state/state_object.go`创建

  ```go
  // internal/ethapi/api.go
  func (s *PersonalAccountAPI) NewAccount(password string) (common.Address, error) {
     ks, err := fetchKeystore(s.am)
     if err != nil {
        return common.Address{}, err
     }
     acc, err := ks.NewAccount(password)
     if err == nil {
        log.Info("Your new key was generated", "address", acc.Address)
        log.Warn("Please backup your key file!", "path", acc.URL.Path)
        log.Warn("Please remember your password!")
        return acc.Address, nil
     }
     return common.Address{}, err
  }
  ```

- **算法生成过程：**

  - **第一步：32字节，私钥 (private key)**

  　　伪随机数产生的256bit私钥示例(256bit)

  　　`18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725`

  - **第二步：64字节，公钥 (public key)**

  　　采用椭圆曲线数字签名算法ECDSA-secp256k1将私钥（32字节）映射成公钥（算上前缀65字节）

  ​	（前缀04+X公钥+Y公钥），公钥是椭圆曲线上的一点，故有`(X, Y)`

  　　`04`<br>
  　　`50863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352`<br>
  　　`2cd470243453a299fa9e77237716103abc11a1df38855ed6f2ee187e9c582ba6`

  ​   （去掉`04`前缀）计算公钥的 **Keccak-256** 哈希值（32bytes）：

  　　`fc12ad814631ba689f7abe67 1016f75c54c607f082ae6b0881fac0abeda21781`

  - **第三步：20字节，地址 (address)**

  　　取上一步hash值的**后20bytes**，加上前缀0x，即以太坊地址：

  　　`0x 1016f75c54c607f082ae6b0881fac0abeda21781`

- **[签名Sign](https://learnblockchain.cn/books/geth/part3/sign-and-valid.html)**（65字节）

![密码学技术分类](https://img.learnblockchain.cn/2019/05/03_cryptography-technology.png!de)

虽然以太坊签名算法也采样了 secp256k1（与比特币相同） ，但是在签名的格式上有所差异，比特币在 [BIP66](https://github.com/bitcoin/bips/blob/master/bip-0066.mediawiki)中对签名数据格式采用严格的DER(Distinguished Encoding Rules，可辨别编码规则)编码格式。

> ECDSA.spec256k1椭圆曲线签名算法
>
> **使用公钥叫加密数据，使用私钥叫签名，签名通过公钥验签**

以太坊的签名格式是`r+s+v`。`r`和`s`是ECDSA签名的原始输出，而末尾的一个字节为恢复id（recovery id简称recid） ，在以太坊中用`v`表示，是签名的最后一个字节。 **65 字节的序列：r 有 32 个字节，s 有 32 个字节，v 有一个字节。**

> recid称为恢复标识符。因为我们使用的是椭圆曲线算法，仅凭 r 和 s 可计算出曲线上的多个点，因此会恢复出两个不同的公钥（及其对应地址）。v 会告诉我们应该使用这些点中的哪一个（也可以理解为查找次数）。在大多数实现中，v =recid，在内部只是 0 或 1。

```go
// crypto/signature_nocgo.go
func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
   // 签名是针对32字节的byte，实际上是对应待签名内容的哈希值
   if len(hash) != 32 {
      return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
   }
   ......
   defer priv.Zero()
  // 调用比特币的签名函数，传入secp256k1 、私钥和签名内容,并说明并非压缩的私钥。
  // 此时 SignCompact 函数返还的签名格式为：[27 + recid] [R] [S]
   sig, err := btc_ecdsa.SignCompact(&priv, hash, false) // ref uncompressed pubkey
   if err != nil {
      return nil, err
   }
  
   // 以太坊签名格式是[R] [S] [V]，和比特币不同。因此需要进行调换位置
   // 减去27的原因是，比特币中第一个字节的值等于27+recid，因此 recid= sig[0]-27
   v := sig[0] - 27
   copy(sig, sig[1:])
   sig[RecoveryIDOffset] = v
   return sig, nil
}
```

> 在以太坊中区块中的数据需要签名的仅有**交易Transaction**
>
> **注意：EIP-155后，在交易签名时，v值不再是recid, 而是 v = recid+ chainID*2+ 35**（旧为v = recid+27)

```go
// core/transaction_signing.go
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
  // 交易签名时，需要提供一个签名器(Signer)和私钥(PrivateKey)。
  // 需要Singer是因为在EIP155修复重放攻击漏洞后，需要保持旧区块链的签名方式不变，
  // 但又需要提供新版本的签名方式。因此通过接口实现新旧签名方式，根据区块高度创建不同的签名器。
  h := s.Hash(tx)
  sig, err := crypto.Sign(h[:], prv)
  if err != nil {
  	return nil, err
  }
  // 将签名结果解析成三段R、S、V，拷贝交易对象并赋值签名结果。最终返回一笔新的已签名交易。
  // 对应前文transaction的结构V、R、S
  return tx.WithSignature(s, sig)
}
```

<div align=center>
<img src="https://img.learnblockchain.cn/2019/04/27_ethereum-tx-sign-flow.png" style="width:65%;">
</div>


#### 2）CA账户

**Storage**是一个 key 和 value 都是`common.Hash`类型的 map 结构，account code 实际存储位置，与EOA账户相比合约账户额外保存了一个存储层(Storage)用于存储合约代码中持久化的变量的数据

> 对应solidity智能合约的状态变量

Storage 层的基本组成单元称为槽 (Slot)。若干个 Slot 按照Stack的方式顺序集合在一起就构造成了 Storage 层。每个 Slot 的大小是 256 bits（32 bytes）的数据。

Slot 作为基本的存储单元，通过地址索引的方式被上层函数访问。Slot的地址索引的长度同样是32 bytes(256 bits)，寻址空间从 `0x0000000000000000000000000000000000000000000000000000000000000000` 到 `0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF`。因此在理论状态下，一个 合约 可以最多保存 2^256 bytes (0 - 2^256-1) 的数据，这是个相当大的数字。



```
# 插槽式的数组存储
----------------------------------
|               0                |     # slot 0
----------------------------------
|               1                |     # slot 1
----------------------------------
|               2                |     # slot 2
----------------------------------
|              ...               |     # ...
----------------------------------
|              ...               |     # 每个插槽大小为 32 字节
----------------------------------
|              ...               |     # ...
----------------------------------
|            2^256-1             |     # slot 2^256-1
----------------------------------
```



> 可观测宇宙中约有2^272个原子

```go
// core/state/state_object.go
type Storage map[common.Hash]common.Hash
```

为了更好的管理数据，Contract 同样使用 MPT(会对`slot[key] = value`中key值(即slot position)进行hash存储) 作为索引树（Account stroage trie）来管理 Storage 层的Slot。

值得注意的是，合约 Storage 层的数据并不会跟随交易一起，被打包进入 Block 中。只有Account storage Trie root 被保存在 account state 结构体中，而account state构成世界状态state。

因此，当某个 Contract 的 Storage 层的数据发生变化时，任一叶子结点account state的改变会使account storage tire root改变，进而使世界状态state其中的一个叶子结点发生改变，进而改变state root，从而记录到Chain链上（包含在block header中）。

> Storage 的数据读取和修改，具体是在执行相关 Transaction 的时候，通过 EVM  opcodes中**sload**和**sstore**来实际执行的
