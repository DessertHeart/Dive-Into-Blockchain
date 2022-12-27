### 区块链Block => Blockchain

#### 1）概念与源码分析

以太坊可以视为一个分布式数据库，数据库数据的变更由交易Transaction驱动。为了有效、有序管理交易，必须将一笔或多笔交易组成一个数据块，才能提交到数据库中。这个数据块即**区块(Block)**。一个区块不但包含了多笔交易，还记录一些额外数据，以便正确提交到数据库中。

```go
// core/types/block.go
type Block struct {
   // 区块的结构构成
   header       *Header
   uncles       []*Header
   transactions Transactions

   // 缓存cache
   hash atomic.Value
   size atomic.Value

   // eth protocol追踪
   ReceivedAt   time.Time
   ReceivedFrom interface{}
}

type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`

	// EIP-1559
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`
}

type Body struct {
	Transactions []*Transaction
	Uncles       []*Header
}
```

**"链"**指的是每个区块加密引用其父块。 换句话说，区块被单向链接在一起。 在不改变所有后续区块的情况下，区块内数据是无法改变，但改变后续区块需要整个网络的共识。

> **分叉fork**：指的是出现了两条互相竞争的链，出现这种不正常的情况后，区块链要依照某种方式选择一条作为canonical chain，即唯一正确的链。**解决 - LMD-GHOST 分叉选择算法**
>
> **重组reorg**：指的是在选择出唯一正确的链之后，又因为某种原因，攻击什么的，这个决定被推翻了，另外一条链又成为了canonical chain。**解决 - Gasper (Casper-FFG)惩罚和罚没**
>
> *<u>注意：以上两个情况均发生于chian header时期</u>*

```go
// core/blockchain.go
type BlockChain struct {
	chainConfig *params.ChainConfig // 链和网络的基本配置
	cacheConfig *CacheConfig        // 节点代码裁剪缓存

	db         ethdb.Database // levelDB物理存储，数据持久化
	snaps      *snapshot.Tree // Snapshot tree for fast trie leaf access
	triegc     *prque.Prque   // Priority queue mapping block numbers to tries to gc
	gcproc     time.Duration  // Accumulates canonical block processing for trie dumping
	triedb     *trie.Database // trieDB,存储trie
	stateCache state.Database // 世界状态State database

	// 可索引追溯得到的最大区块范围
	//  * 0:   means no limit and regenerate any missing indexes
	//  * N:   means N block limit [HEAD-N+1, HEAD] and delete extra indexes
	//  * nil: disable tx reindexer/deleter, but still index new blocks
	txLookupLimit uint64

  // chain header：负责特殊的数据维护，每次加入一个新的block都会判断
  // validator们根据自己当前看到的场景，来判断哪个block是beacon chain头的过程
  // 包括: (1) total difficulty (2) header (3) block hash -> number mapping 
  // (4) canonical number -> hash mapping (5) head header flag.
	hc            *HeaderChain
  
  // Feed: channel实现的一对多的订阅式消息分发结构
	rmLogsFeed    event.Feed
	chainFeed     event.Feed
	chainSideFeed event.Feed
	chainHeadFeed event.Feed
	logsFeed      event.Feed
	blockProcFeed event.Feed
	scope         event.SubscriptionScope
	genesisBlock  *types.Block
	......
}
```

#### 2）[共识机制](https://ethos.dev/beacon-chain)

- **共识逻辑**

  slot是将块添加到信标链的时机，每个slot为 12 秒，一个epoch为 32 个slot：6.4 分钟

  ![Checkpoint is same for Epoch 1 and Epoch 2](https://ethos.dev/assets/images/posts/beacon-chain/Beacon-Chain-Checkpoints.jpg)

  > 注意：某个 slot 可能没有区块，但当系统以最佳方式运行时，区块会添加到每个可用的 slot 中
  >
  > Beacon Chain genesis block 位于 Slot 0

  每个 epoch 中，第一个 slot 中的区块是一个**检查点 (checkpoint)**。被用于使区块链账本上的记录变得永久和不可篡改。

  - 第一步，如果所有活跃验证者质押的 ETH 余额中至少有 2/3 (即“绝对多数”) 证明了**最近的两个检查点**(当前的被称为“目标检查点”，前一个被称为“源检查点”)，那么这两个检查点之间的这段区块就被认定为**“合理化” (justified) **。

  - 第二步，一旦某个被“合理化”的检查点之后新出现了另一个被“合理化”的检查点，那么前一个检查点就是**“敲定” (finalized)** 。在这个检查点之前的所有区块/记录都成为了区块链上永久不可篡改的记录。

- **系统角色**

  - Validator验证者：虚拟矿工Miner，验证从其他验证者那里接收到的区块并验证正确性并签名，同时对有效的区块进行“**证明”(attestations-包含 LMD GHOST 投票和 FFG 投票)**，并通过**Casper FFG（finality gadget）投票**（FFG投票每个验证者都要做）来对每个epoch的记录进行上述 “justified” 与 "finalized"。此外，任一验证者将不定期地被要求成为**提议者 (proposer)出一个新区块**

    > 注意：每个slot的proposer选取是从j合中选，与committee无关，且还要避免是committee成员。

    

    ![validator key schematic](https://ethereum.org/static/8b68c1825d524f8102b5e58574824c77/e0885/validator-key-schematic.png)

    > staking赛道解决方案：分布式验证者 (Distributed Validators, DV) 是一种将单个验证者/节点的职责被分给几个共同验证者/节点 (co-validator)，应用[BLS 签名方案](https://www.ethereum.cn/Eth2/distributed-validator-specs)以提高与在一个单一机器上运行一个验证者客户端相比的韧性 (安全性、活性，或两者兼有)。如[Obol Network](https://mirror.xyz/bitcoinorange.eth/gyXAG1neqkm7nBCNFk77wLd5llZIFoVk3eqezI-wPZI)
    >
    > 很多validator都会跑一个backup，因为validator不能下线。如果主程序和backup都跑的话，是会出现发送两个attestation的这个情况（如一个epoch发送两个attestation），该情况是没有办法从protocol层面限制的，因为attestation发送出去，是让不同的人收集的，只要验证者签名合法的，就会被其他人收集起来，其他人没有能力判断发送的attestation是否重复

  - Committees委员会 (验证者子集) ：**同步委员会 (sync committee)** 由随机分配**（RANDAO）**至少128名验证者组成的小组，一个验证者在一个epoch内只能在一个委员会中，每个epoch里，各个委员会被均分给每个slot（1slot - 1 committee），每27小时更新一次。

    这些被随机选中的验证者，除了其验证者本职工作以外，还将对**Chain head**进行签名（如遇分叉情况，依据LMD GHOST分叉选择算法，于chain header timing时期）。

    > 轻客户端可以检索这些被验证过的区块，而无需访问整条历史链或整个验证者集。

- **[奖惩机制](https://eth2book.info/altair/part2/incentives/rewards)**

  - 奖励：验证者以 32 ETH 用于作为“抵押品”。

    - 当验证者进行的 **LMD-GHOST 投票**和 **FFG 投票**与大多数其他验证者一致时，那么验证者就会获得证明奖励。

    - 当验证者被选中作为“**区块提议者”(block proposer) **时，如果其提议的区块被“敲定”， 那么该验证者也将获得奖励。此外，区块提议者也可以通过将有关其他验证者行为不当的证明打包进自己提议的区块中，从而增加自己获得的奖励（“报酬”）。

      > block proposer的reward中包括交易的gas priority fee（execution layer），还有consensus layer中专门有一个incentive layer，proposer出块会有奖励（类似对应PoW时的miner挖出区块奖励，PoS之后的出块奖励变到这里，**数量是其区块中包含总奖励的1/7**）
      >
      > rewards在执行交易的同时就去分发了，beacon chain有维护一个validator的balance的mapping，所有的奖励都在这个map里，每个epoch重新分配委员会时，会根据将本地validator的balance更新上链，但实际真正有效的balance是直到此次epoch为止，最新finalized区块的位置对应validator的balance，大家也都会去查这个时候的余额，因为是确定有效且不可篡改的

  - 惩罚

    - 惩罚Penalties：以各种机制的形式来销毁一部分验证者质押的 ETH

      - 验证者未能提交一个 FFG 投票、提交延迟了或者提交了错误的 FFG 投票时，会受到**证明惩罚 (attestation penalties)**。削减的数额等同于其提交正确的证明而原本可以获得的奖励。
      - 验证者错过了进行 LMD-GHOST 投票，则不会受到惩罚，只是错过了本可以通过对链头进行投票而获得的奖励。

    - 罚没Slashings：验证者发生严重行为(如下列举)，会导致验证者被强制从网络中**移除**，其质押的 ETH 的 1/64 (最高可达 0.5 ETH) 将立即被销毁，然后开始一个**为期 36 天的移除期**，在此期间，验证者的质押金将逐渐被削减；且在这段期间的中间点 (**第18天**)，该验证者还将受到额外的惩罚，惩罚大小将与此次罚没事件发生之前的 **36 天内**所有被罚没的验证者的 ETH 质押总额成比例（**串谋惩罚correlation penalty**）

      - 在同一个 slot 提议和签名两个不同的区块。
      - 对“环绕”某个区块的另一个区块进行证明 (实际上就是更改区块链历史)。
      - 通过对同一个区块的两个候选区块进行“双重投票”(double voting)。

      > 1、对于节点来说，无论有意无意，如果有对应行为上链并触发slashing，那就是无法补救的了
      >
      > 2、protocol层和程序层都没有办法限制恶意节点多发attestation，有的方式只有合法节点的打包“举报”（事实上会有很多节点在盯着，因为有奖励）。但除了slashing罚款，只要所谓的“声明无效证明”，但已经造成的链上既定事实是没办法改变的。所以，对于恶意节点来说，成本只有32ETH，但成功的几率权衡巨大的获利，还是可能预见的，所以很多攻击者的做法就是集中发送很多相互冲突的attestation

    - Inactive Leak机制：如果信标链已经有超过 4 个 epoch 都没有被敲定时触发（超过 1/3 的验证者离线或未能提交证明的证明，即不可能有 2/3 的绝对多数验证者来敲定检查点）。**逐渐削减不活跃的验证者的 ETH 质押金，直到这些验证者控制的质押金少于网络中总质押金的 1/3，从而允许剩余的活跃验证者对区块链进行敲定**。无论这些不活跃的验证者数量有多大，剩余的验证者最终都将控制 >2/3 的总质押金。
