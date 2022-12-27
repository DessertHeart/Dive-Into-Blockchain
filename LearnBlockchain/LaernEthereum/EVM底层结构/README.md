## EVM底层结构

> 以太坊是一个状态机，最终区块交易的状态就是区块链的状态



### 1、数据结构

![以太坊区块结构](https://img.learnblockchain.cn/2019/05/19_ethereum-full-block-data-struct.png!de)

#### 1) accountSate帐户

| accountSate                                                  | 描述说明                                                     |      |
| :----------------------------------------------------------- | :----------------------------------------------------------- | ---- |
| [nonce](https://www.8btc.com/books/834/ethereum-book/_book/%E7%AC%AC%E4%B8%83%E7%AB%A0.html) | 基于0的计数器，账户发起交易计数，对于合约账户，表示此帐户创建的合约数量并初始化为1 |      |
| balance                                                      | 余额                                                         |      |
| storageRoot                                                  | 账户存储树根哈希                                             |      |
| codeHash                                                     | 如果有关联代码就是合约账户                                   |      |

#### 2) transaction交易

| transaction | 描述说明                   |
| :---------- | :------------------------- |
| nonce       | 发出交易的账户的nonce值    |
| gasPrice    | 余额                       |
| gasLimit    | transaction里最多能用费用  |
| to          | 接受者地址                 |
| value       | 转账金额、合约创建初始充值 |
| data        | 附加数据                   |
| v,r,s       | 交易签名                   |

#### 3) receipt收据

| receipt           | 描述说明                                                |
| :---------------- | :------------------------------------------------------ |
| postState         | []byte，交易之后的状态的RLP编码后结果                   |
| cumulativeGasUsed | 交易总共gas消耗余额                                     |
| logsBloom         | [256]byte，布隆过滤器，存储了一笔交易的收据中所有的日志 |
| logs              | []*Log，存储涉及的日志                                  |

#### 4) block区块（Header + Body）

| block Header     | 描述说明                                                     |
| :--------------- | :----------------------------------------------------------- |
| **parentHash**   | **父块头哈希，正是通过此记录，才能完整的将区块有序组织，形成一条区块链** |
| sha3Uncles       | 叔块头哈希列表的哈希（PoS不存在该概念,uncles=`[]`，该值固定为`0xc0=RLP([])`） |
| beneficiary      | 挖矿奖励获取人地址address（miner）                           |
| stateRoot        | 状态树根节点Hash                                             |
| transactionsRoot | 交易树根节点Hash                                             |
| receiptsRoot     | 收据树根节点Hash                                             |
| logsBloom        | [256]byte，存储了当前区块中所有的receipt收据的日志的布隆过滤器 |
| difficulty       | big.Int，PoW时难度值 12s-16s自动调整（PoS后固定为0）         |
| number           | big.Int，区块高度                                            |
| gasLimit         | uint64，区块内能用的gas 允许矿工可以有5%的上下浮动（伦敦升级后为1500W，但块的大小将根据网络需求增加或减少，直到块限制为 3000 W{原目标大小的 2 倍}，solidity合约大小限制为的 `24576bytes`） |
| gasUsed          | uint64，区块内所有交易实际消耗的gas                          |
| timeStamp        | uint64，表示此区块创建的UTC时间戳，单位秒（PoW矿工可以将区块时间戳修改 +/-15 秒，[PoS为固定12s，极少数会有12的整数倍的情况，会检查时间戳](https://ethereum.stackexchange.com/questions/135445/miner-modifiability-of-block-timestamp-after-the-merge)） |
| extraData        | <=32字节数组，由矿工自定义，一般会写一些公开推广类内容或者作为投票使用。 |
| **mixHash**      | 本区块标识哈希（hash(区块头数据不包含nonce)）                |
| nonce            | uint64（8字节）随机数，PoW工作量证明（PoS后固定为0）         |

| block Body      | 描述说明                          |
| :-------------- | :-------------------------------- |
| TransactionList | 交易列表                          |
| ommersList      | 叔块头哈希列表(PoS后固定为空`[]`) |

#### 5）四棵全局树（Trie）

- ##### World State Trie世界状态树

| Key                        | value             | 描述说明                    |
| :------------------------- | :---------------- | :-------------------------- |
| keccak256(ethereumAddress) | RLP(accountState) | 关联block header的stateRoot |

- ##### Transactions Trie区块级交易树

| key                   | value            | 描述说明                           |
| :-------------------- | :--------------- | :--------------------------------- |
| RLP(transactionIndex) | RLP(transaction) | 关联block header的transactionsRoot |

- ##### Receipts Trie区块级数据树

| Key                   | value                   | 描述说明                       |
| :-------------------- | :---------------------- | :----------------------------- |
| RLP(transactionIndex) | RLP(transactionReceipt) | 关联block header的receiptsRoot |

- ##### Account Storage Trie账户存储树

| key                      | value              | 描述说明                     |
| :----------------------- | :----------------- | :--------------------------- |
| keccak256(slot position) | RLP(slot position) | 关联acountState的storageRoot |

*<u>注意：四个Trie 都不在区块链网络传输 ，网络中只传输交易数据</u>*



### 2、存储设计

**以太坊虚拟机（或[EVM](https://ethereum.org/en/developers/docs/evm/)）是基于堆栈的计算机**。这意味着所有指令都从堆栈中获取它们的参数，并将它们的结果写入堆栈。因此，每条指令都有堆栈输入、它需要的参数（如果有的话）、堆栈输出和返回值（如果有的话）。所有指令都编码为 1 个字节，PUSH 指令除外，它允许将任意值放入堆栈并在指令后直接对该值进行编码。可用指令列表及其操作码显示在[参考资料中](https://www.evm.codes/)。

> [虚拟机VM](https://zhuanlan.zhihu.com/p/53692225)：运行在真实机器上的软件，提供操作系统（在系统VM的情况下）或应用程序（在进程 VM的情况下）的运行环境。

> 传统CPU以及诸如 Dalvik 虚拟机（Android移动设备平台的核心组成部分之一，可以支持已转换为.dex 格式的Java「注：JVM是基于堆栈的」应用程序的运行），是基于寄存器的结构：[区别](https://yanyezhang.github.io/2018/08/15/%E3%80%90Java%E8%99%9A%E6%8B%9F%E6%9C%BA%E7%AE%80%E5%8F%B2%E3%80%91%E5%9F%BA%E4%BA%8E%E6%A0%88%E8%99%9A%E6%8B%9F%E6%9C%BAvs%E5%9F%BA%E4%BA%8E%E5%AF%84%E5%AD%98%E5%99%A8%E8%99%9A%E6%8B%9F%E6%9C%BA/)

| 存储结构 |                             功能                             |                          规则与大小                          |
| :------- | :----------------------------------------------------------: | :----------------------------------------------------------: |
| ROM      | 用来保存所有EVM程序代码的“只读”存储，由以太坊客户端独立维护  |                  存储只读代码，Code is Law                   |
| Stack    |      即所谓的“运行栈”，用来保存EVM指令的输入和输出数据       |   256 bit(32字节)位宽，最大深度为1024   |
| Momory   | 内存，一个简单的字节数组，用于临时存储EVM代码运行中需要的存取的各种数据 |     256bit / 8bit位宽，无限大小，基于32字节进行寻址和扩展      |
| Storage  |         存储，由以太坊客户端独立维护的持久化数据区域         | 256 * 2 bit位宽(key & value)，2^256大小，每个账户的存储区域被以32字节为单位划分为若干”槽（slot）“，合约中的“状态变量”会根据其具体类型分别保存到这些“槽”中 |

#### 1)  storage, memory, calldata, stack区分

在 Solidity 中，有两个地方可以存储变量 ：**存储（storage）以及内存（memory）**。Storage变量是指永久存储在区块链中的变量。Memory 变量则是临时的，只能用于函数内部，当外部函数对某合约调用完成时，内存型变量即被移除。

>  storage 在区块链中是用key/value的形式存储slot[x]=y，而memory则表现为字节数组[0x01, 0x02...]

内存(memory)位置还包含2种类型的存储数据位置，一种是calldata，一种是栈（stack）

**(1) Calldata**这是一块**只读**的，且不会永久存储的位置，用来存储函数参数。 外部函数的参数（非返回参数）的数据位置被强制指定为 calldata ，效果跟 memory 差不多。 

**(2) Stack** ，EVM是一个基于栈的语言，栈实际是在计算机内存中的一个数据结构，每个栈元素占为256位，栈最大长度为1024。 **值类型的局部变量是存储在栈上，但注意，堆栈仅有高处的 16 层是可以被快速访问的（solidity中stack overflow报错的原因）**

> 操作码Opcodes（字节指令）：EVM指令被分配了一个介于 0 和 255（或十六进制中的 FF）之间的值。它是帮助我们人类阅读指令的文本表示。智能合约是一组指令。当 EVM 执行智能合约时，它会逐条读取并执行每条指令。如果无法执行指令（例如，因为堆栈上没有足够的值），则执行将返回。

![../_images/Picture48.png](https://ethbook.abyteahead.com/_images/Picture48.png)

### 3、编码方式

#### 1）[HP编码](https://www.jianshu.com/p/8b6d2a7fb6b5)

十六进制前缀 Hex-Prefix 编码：

> 主要树节点结构 Merkle Patricia tree （Trie）

#### 2）[RLP编码](https://learnblockchain.cn/books/geth/part3/rlp.html)

递归长度前缀编码 Recursive Length Prefix Encoding：

> EVM只认序列化数据，state数据保存与传输，将任意格式的数据编码串型
>
> BE 是将正整数值扩展为最小长度的高端字节数组的函数，点运算符是执行序列拼接

![以太坊技术与实现-图-以太坊RLP 编码算法-数据标记规则](https://img.learnblockchain.cn/book_geth/2019-12-28-23-20-21.png!de?width=700px)

#### 3）RLP与[ABI-Encoding](https://me.tryblockchain.org/Solidity-abi-abstraction.html#toc_5)

应用二进制接口ABI（ Application Binary Interface）是与合约交互的标准。 EVM使用 ABI 编码的数据来理解要执行字节码的哪一部分。合约交互只是以太上的一种交易。payload（要做什么） 位于transaction的`data`字段中。

调用合约的用户通过 ABI-Encoding 对的输入参数和函数签名进行编码（当然也包括合约代码的输出参数，比如日志Log的data），并将其放在transaction的`data`字段中。之后队transaction进行签名作为所有权证明。最后，对整个transaction数据进行 RLP 编码，即可传递给EVM。

> 注意：ABI编码仅用于合约交互，transaction的`data`字段包含目标函数选择器和函数所需入参。而 RLP，可以理解为是Ethereum在数据通过传输之前进行数据编码(转换)的低级方法

具体合约交互的数据转换逻辑：

- 发送方：
  1. 确定在合约地址与对应的函数
  2. ABI 编码函数选择器和输入参数
  3. 将 ABI 编码器的输出原始字节放入交易的`data`
  4. 将所有值放入交易包括 to、 value 和 nonce
  5. RLP 对这些字段进行编码
  6. 签名交易
  7. 再次RLP 编码并发送
- 接收方：
  1. RLP解码收到的transaction
  2. 验证签名，检查交易的有效性
  3. 再次RLP解码
  4. 执行 (EVM 将在内部处理 ABI 解码与执行)

### 4、EVM细节

- EVM 代码执行的实际gas消耗与其对内存memory的使用有关，并不是固定的。
- 鼓励最小化使用存储storage，用sstore操作将非0值存储区域重置为0值，会获得实时的gas返还。
- 交易执行的最后会删除执行过程中接触过的所有“空账户”和自毁列表中的账户，这也会返还一定量的gas。
- **EVM代码的执行必定会持续到一个正常终止或一个异常终止，但无法用代码直接触发一个异常终止**。
- **EVM代码执行的异常终止会撤销当前交易中所有对状态的更改，但执行过程中所有消耗的gas不会返还**。

