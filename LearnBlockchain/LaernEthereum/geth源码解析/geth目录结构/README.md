### [geth目录结构](**https://github.com/DessertHeart/go-ethereum**)(2022.11 master)

> + **accounts**：账户管理
>
> + **Build**：makefile编译与构建相关脚本文件
> + **cmd**：命令行工具管理
>   + abidump：解析给定ABI并尝试从4byte数据库解释
>   + **abigen**：源码生成器，将solidity合约 => safe-go包（golang友好工具）
>   + **Bootnode**：启动网络发现节点，bootnode是指以太坊p2p通讯时的初始引导节点，在找不到任何peers时，就去连接bootnode（通过一系列bootstrap引导程序实现）
>   + checkpoint-admin：更新检查点oracle状态的工具（如：部署检查点oracle合约、签署新的检查点以及更新检查点oracle契约中的检查点）
>   + clef：签署交易和数据的工具，帮助dapp不依赖Geth的账户管理
>   + devp2p：测试p2p通信工具
>   + ethkey：处理以太坊kestore密钥文件的工具
>   + **evm**：虚拟机开发工具，受隔离的代码调试环境，可以用于调试opcodes
>   + faucet：水龙头工具
>   + **geth**：以太坊命令行客户端
>   + p2psim：提供工具以模拟http的API
>   + puppeth：新以太坊网络创建向导工具
>   + rlpdump：RLP序列化输出工具
>   + Utils：提供一些公共工具
> + **common**：公共工具包管理
> + **consensus**：共识算法（[beacon-Gasper](https://mp.weixin.qq.com/s/3Gw3DuBr-LCm0k9NaAu3bg)（2.0-POS），[ethash](https://learnblockchain.cn/books/geth/part2/consensus/ethash.html)（1.0-POW），[clique](http://yangzhe.me/2019/02/01/ethereum-clique/)（private testnet-POA））
> + **console**：geth console交互式命令管理（Geth Web3 控制台的入口）
> + **contracts/checkpointoracle**：[检查点checkpoint oracle](https://learnblockchain.cn/article/901)以太坊预言机的合约实现与管理
> + **core**：核心数据结构（区块，EVM，logsbloom...）
> + **crypto**：密码学相关算法
> + **docs**：代码审计与一些事件文件管理
> + **eth**：以太坊协议层的实现
> + **ethclient**：RPC客户端（geth）
> + **ethdb**：数据库管理（世界状态stateDB底层levelDB以及一些测试用数据库）
> + **ethstats**：以太坊网络状态报告服务
> + **event**：实时事件管理
> + **graphql**：query language，针对graph（图状数据）高效查询语言
> + **internal**：一些内部使用的工具库的集合
> + **les**：以太坊轻量级轻节点通讯子协议（LES），为轻节点通讯准备的
> + **light**：以太坊轻客户端相关功能实现，可按需检索的状态和链对象等
> + **log**：节点运行的日志（人机友好格式）管理
> + **metrices**：服务监控相关代码
> + **miner**：区块创建于挖矿相关逻辑
> + **mobile**：geth移动端相关封装与api
> + **node**：以太坊各种节点的实现
> + **p2p**：p2p网络协议
> + **params**：预设一系列参数值（如: bootnode的enode地址）
> + **rlp**：实现递归长度前缀，以太坊序列化编码方法
> + **rpc**：RPC接口封装（如IPC、http、websocket）
> + **signer**：签名signature的实现与管理
> + [**swarm**](https://www.ethswarm.org/)：swarm分布式存储实现
> + **tests**：单元测试相关代码
> + **trie**：MPT（Merkle Patrica Trie）数据结构的实现

