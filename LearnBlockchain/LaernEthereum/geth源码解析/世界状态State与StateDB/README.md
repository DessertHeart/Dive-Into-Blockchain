### 世界状态Ethereum State与StateDB

#### 1）以太坊状态机（[接收器/识别器型状态机](https://learnblockchain.cn/books/geth/part3/statedb.html)）

`有限自动机` （Finite Automata Machine）是计算机科学的重要基石，它在软件开发领域内通常称作`有限状态机`（ Finite State Machine，缩写 FSM），简称**状态机**，是表示有限个状态以及在这些状态之间的转移和动作等行为的数学模型。

​																		$次态 = f(现态，输入)$

通过节点维护运行，以太坊网络是一个去中心化状态机。在任意时刻，只会处于一个全世界唯一的状态，我们把这个状态机，称之为以太坊世界状态，代表着以太坊网络的全局状态。

**世界状态(state)**由无数的账户信息组成，每个账户均存在一个唯一的账户信息。账户信息中存储着账户余额、Nonce、合约哈希、账户状态等内容，每个账户信息通过账户地址影射。 从创世状态开始，随着将交易作为输入信息，在预设协议标准（条件）下将世界态推进到下一个新的状态中。

![以太坊技术与实现-状态](https://img.learnblockchain.cn/book_geth/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E6%8A%80%E6%9C%AF%E4%B8%8E%E5%AE%9E%E7%8E%B0-%E5%9B%BE2019-12-7-23-35-20!de?width=600px)

#### 2）[StateDB](https://learnblockchain.cn/books/geth/part3/statedb.html)

StateDB是EVM State中最高层的封装，直接提供了与StateObject (Account，Contract)相关的 CURD 的接口给其他的模块，充当状态（数据）、Trie(树)、LevelDB（存储）的协调者。

- 从程序设计角度，StateDB 有多种用途：
  1. 维护账户状态到世界状态的映射。
  2. 支持修改、回滚、提交状态。
  3. 支持持久化状态到数据库中。
  4. 是状态进出默克尔树的媒介。

> 需要注意，世界状态中的所有状态都是以`Account`账户为基础单位存在的。所访问的任何数据必然属于某个账户下的状态，世界状态仅仅是通过一颗树来建立安全的映射。比如所访问的数据可以分为如下几种类型：
>
> 1. 访问账户基础属性：`Balance、Nonce、Root、CodeHash`
> 2. 读取合约账户代码
> 3. 读取合约账户中存储内容
>
> 在代码实现中，为了便于账户隔离管理，使用不开放的 `stateObject` (见账户结构)来维护

1、`trie root`：首先，我们要告诉 StateDB ，我们要使用哪个状态。因此需要提供 StateRoot 作为默克尔树根去构建树。StateRoot 值相当于数据版本号，根据版本号可以明确的知道要使用使用哪个版本的状态。

2、`database`：然后，数据内容本身并没在树中，需要到具体数据库中读取。因此在构建 StateDB 时需要提供 state root 和 db 才能完成构建。

> 轻节点使用的 odrDatabase，对数据读取方式的封装，因为需要通过向其他节点查询来获得数据

```go
// core/state/database.go
type Database interface {
  // 打开指定状态版本(root)的含state trie
	OpenTrie(root common.Hash) (Trie, error)
  // 打开账户(addrHash)下指定状态版本(root)的Account storage trie。
	OpenStorageTrie(addrHash, root common.Hash) (Trie, error)
  // 深度拷贝树
	CopyTrie(Trie) Trie
  // 获取账户（addrHash）的合约，必须和合约哈希`codeHash`匹配
	ContractCode(addrHash, codeHash common.Hash) ([]byte, error)
  // 获取指定合约大小
	ContractCodeSize(addrHash, codeHash common.Hash) (int, error)
	// 获得 Trie 底层的数据驱动 DB，如: levedDB 、内存数据库、远程数据库
	TrieDB() *trie.Database
}
```

然后，即可初始化一个stateDB

```go
// core/state/statedb.go
func New(root common.Hash, db Database, snaps *snapshot.Tree) (*StateDB, error) {
  // 1.trie: 打开指定状态版本(root)的含世界状态的顶层树
  tr, err := db.OpenTrie(root)
	......
  // 2.初始化stateDB
	sdb := &StateDB{
    // key point1： database
		db:                  db,
    // key point2： trie
		trie:                tr,
		......
	}
	......
}
```

- **持久化**

> ​	`journal`参数(from struct statedb)： 记录修改状态的日志流水，使用此日志流水可回滚状态

在区块中，**交易**作为输入条件，来根据一系列动作修改状态。 

-  **StateDB 可视为一个内存数据库**，在完成区块挖矿前，只是获得在内存中的状态树的 Root 值。状态数据先在内存数据库中完成修改，所有关于状态的计算都在内存中完成。 

-  **在将区块持久化时，完成有内存到数据库（真正落盘）的更新存储**，此更新属于增量更新，仅仅修改涉及到被修改部分。

```go
// core/state/statedb.go
func (s *StateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	......
}

```

![以太坊技术与实现-图以太坊 State 库读写关系](https://img.learnblockchain.cn/book_geth/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E6%8A%80%E6%9C%AF%E4%B8%8E%E5%AE%9E%E7%8E%B0-%E5%9B%BE2019-12-18-21-56-7!de?width=600px)

> 底层物理存储层DB只有 LevelDB，为了提高读写性能，使用 cachingDB 对其进行一次封装，使用了LRU缓存淘汰算法。
>
> [LevelDB](https://learnblockchain.cn/article/728)：持久化KV单机数据库，具有很高的随机写，顺序读/写性能，但是随机读的性能很一般，也就是说，**LevelDB很适合应用在查询较少，而写很多的场景**。LevelDB应用了LSM (Log Structured Merge) 策略，lsm_tree对索引变更进行延迟及批量处理，并通过一种类似于归并排序的方式高效地将更新迁移到磁盘，降低索引插入开销。

#### 3）Merkle Patricia Trie

> Trie 是一种有序的树结构，用于存储和检索键值对（key-value），其中 key 可以映射到有限“字符集”组成的字符串，树的每个节点记录了一个字符，并且指向了下一个字符，每个路径可以组成一个完整的 key，这使得节点可以**共享相同的前缀**。
>
> “Trie” 一词提取自 “re**trie**val”（数据检索）的中间部分，根据其特征，也叫前缀树（Prefix Tree）、字典树等

MPT = Merkle Tree(节点存储数据块的哈希) + Patricia Trie（压缩前缀树，以节省空间高效查询，如图）

![img](https://img-blog.csdnimg.cn/0f01c681c7454bc693bb55e592a4de77.png)

- **MPT节点类型：**

  - 空白节点 NULL

  - 分支节点 branch [ v0 ... v15, vt ]长度为 17 的数组，前 16 个元素表示十六进制字符集，最后一个元素存储该分支对应的value（如果存在）

    > 减小了每个分支节点的容量，但是在一定程度上增加了树高。

  - 叶子节点 leaf [encodedPath, value]

  - 拓展节点 extension [encodedPath, key]

- **以太坊MPT中Key的定义**

  1、key的存储内容（两种）：

  - Origin Key：数据的原始 key，为字节数组（RLP编码）。
  - Secure Key：为原始 key 计算哈希 Keccak256(Origin Key) 的结果，长度固定为 32 字节，用于防止深度攻击。后文我们将看到**以太坊的状态树和存储树**使用这种 Key 类型

  2、key的存储形式：

  - Hex Key：将 Origin Key 或 Secure Key 进行半字节（nibble）拆解后的 key，为 **MPT 树真正存储的 key**。在以上条件的限制下，MPT 树 key 的长度固定为 64 字符（32字节对应）。

    其中的一个必要优化手段是HP Key：hex prefix encoding，Hex 前缀编码。当我们使用 nibble 寻找路径时，我们可能最后会剩下奇数个的 nibble，但是由于数据存储的最小单位是字节，所以可能会带来一些歧义，比如我们可能无法区分 1 或 01（都存储为1字节`01`）。因此，为了区分奇偶长度，叶子节点和拓展节点的 encodedPath 使用一个前缀作为标签，另外，这个标签也用于区分节点类型。

> **nibble**：占4bits ，一位十六进制数即半字节。为HP编码中hex用到的数据结构单位，可以表示数字 0~15，这一步可以看成是将 key 映射到十六进制字符 0~f 组成的字符串，这就是为什么分支节点的数组长度为 17（16+1）

![img](https://img-blog.csdnimg.cn/888b0ba0ef994b7cad561b43262b3a62.png)

```go
// trie/trie.go
type Trie struct {
	root  node
	owner common.Hash

	// 记录从上次哈希操作至今，插入叶子结点leaves叶数量
	unhashed int

	// 检索trie各节点的handler trie工具
	reader *trieReader

	// tracing trie变更的工具, 一个调用合约的交易在执行过程中，可能会改变很多state variable，
  // 它每一步具体都改变了什么，都在trace中记录
	// 每次commit操作会重置
	tracer *tracer
}
```

应实现方法：

```go
// core/state/database.go
type Trie interface {
	GetKey([]byte) []byte
	TryGet(key []byte) ([]byte, error)
	TryGetAccount(key []byte) (*types.StateAccount, error)
	TryUpdate(key, value []byte) error
	TryUpdateAccount(key []byte, account *types.StateAccount) error
	TryDelete(key []byte) error
	TryDeleteAccount(key []byte) error
	Hash() common.Hash
	Commit(collectLeaf bool) (common.Hash, *trie.NodeSet, error)
	NodeIterator(startKey []byte) trie.NodeIterator
	Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error
}
```

### 
