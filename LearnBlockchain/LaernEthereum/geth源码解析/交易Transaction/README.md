

### 交易Transaction

![以太坊交易数据结构](https://img.learnblockchain.cn/2019/04/24_transaction-struct.png!de)

#### 1）Transaction的定义

Transaction是驱动Ethereum执行数据操作的媒介，它主要起到下面的几个作用:

1. 在以太坊网络上的Account之间进行Native Token的转账。
2. 创建新的Contract。
3. 调用Contract中会修改目标Contract中持久化数据或者间接修改其他Account/Contract数据的函数。

> 一、Transaction只能创建Contract账户，而不能用于创建外部账户(EOA)。
>
> 二、如果调用Contract中的只读函数，是不需要构造Transaction的。相对的，所有参与Account/Contract数据修改的操作都需要通过Transaction来进行。
>
> 三、广义的Transaction只能由外部账户(EOA)构建。Contract是没有办法显式构造交易的，但是，实际在合约函数的执行过程中，Contract可以通过构造**internal transaction（合约内部交易称为internal transaction，即在一笔交易执行过程中，合约根据一定条件，进行转账或者是调用新合约等一系列动作产生的结果，正如 etherscan 上标注的一样，`Internal Transactions as a result of Contract Execution`）**来与其他的合约进行交互。

```go
// core/types/transaction.go
type Transaction struct {
	// TxData为接口
  // 与Transaction直接相关的数据都由这个变量来维护
  inner TxData   
	time  time.Time 

  // 缓存：对一些哈希运算等进行缓存，降低CPU计算量
  // 交易哈希
	hash atomic.Value
  // 交易大小：交易信息进行RLP编码后的数据大小,代表交易网络传输大小、代表交易占区块大小、代表交易存储大小。   // 每笔交易进入交易池都需要检查交易大小是否超过32KB
	size atomic.Value
  // 交易发送方EOA
	from atomic.Value
}
```

根据[EIP-2718](https://eips.ethereum.org/EIPS/eip-2718)的设计，原来的`TxData`现在被声明成了一个interface，而不是定义了具体的结构。以便于后续版本对Transaction类型进行更加灵活的修改。

目前，Ethereum中按照时间上的定义了三种类型实现`TxData`接口的的Transaction类型：`LegacyT`，`AccessListTx`，`TxDynamicFeeTx`。

- `LegacyTx`是原始的Ethereum的Transaction设计，目前市面上大部分早年关于Ethereum Transaction结构的文档实际上都是在描述`LegacyTx`的结构（如上transaction结构示意图）。

- `AccessListTX`是基于[EIP-2930](https://support.token.im/hc/zh-cn/articles/900004927906-%E4%BB%A5%E5%A4%AA%E5%9D%8A%E6%9F%8F%E6%9E%97-Berlin-%E5%8D%87%E7%BA%A7%E5%85%AC%E5%91%8A)(柏林分叉)的Transaction。

- `DynamicFeeTx`是[EIP-1559](https://eips.ethereum.org/EIPS/eip-1559)(伦敦分叉)生效之后的默认的Transaction。

#### 2）Transaction的执行流程(事务型)

Transaction的执行主要发生在两个Workflow中:

1. **Miner在打包新的Block时**。Miner首先从从memPool内存池中拿出若干的transaction，按照Gas Price和Nonce对交易进行排序，然后按该顺序对Block中Transaction进行commit执行。`miner/worker.go => mainLoop()`

   > 因为使用mempool有被[抢跑Front-Running](https://github.com/AmazingAng/WTF-Solidity/tree/main/S11_Frontrun)的风险。在解决方案上，有一种使用暗池的方式，即用户发出的交易将不进入公开的`mempool`，而是直接到矿工手里，例如 flashbots 和 TaiChi。也有关于“分歧终端机”[commit-reveal scheme](https://www.geekmeta.com/article/1239071.html)的研究

2. **其他节点添加Block到Blockchain时**。当节点从网络中监听并获取到新的Block时，它们会执行Block中的Transaction，来更新本地的State Trie的 Root，并与Block Header中的State Trie Root进行比较，来验证Block的合法性。`core/blockchain.go => InsertChain()`

一条Transaction执行，可能会涉及到多个Account/Contract的值的变化，最终造成一个或多个Account的State的发生转移。在Byzantium分叉之前的Geth版本中，在每个Transaction执行之后，都会计算一个当前的State Trie Root，并写入到对应的Transaction Receipt中，但这同时会带来大量的冗余计算。因此在**Byzantium分叉**（**2017 年 10 月 16 日**）之后，在一个Block的验证周期中只会计算一次的State Root，一个Block中所有Transaction执行的结果累计使得World State发生一次性的状态转移(具体通过trace记录)。

### 
