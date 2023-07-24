### 以太坊初始化

#### 前言）以太坊客户端与团体


<div align=center>
<img src="https://mmbiz.qpic.cn/sz_mmbiz_png/tsoSYGv5wmO9lNWiaMNmzUApfBFpUnN2nz5ibQsWMwiaxClPEEHum1nSxmlmgYdtLhAUD95oPiaAPjAMasHWibYxSUA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1" style="width:65%;">
</div>

*<u>2022年1月源自[ethernodes](https://ethernodes.org/)与[blockprint](https://github.com/sigp/blockprint)</u>*

- *以太坊客户端*

<div align=center>
<img src="https://ethereum.org/static/7a59bdb7a666b01a74535e8bea21a532/c1b63/eth1eth2client.png" style="width:65%;">
</div>

  - **执行客户端**(原ETH1.0客户端)：也称为执行引擎，侦听网络中广播的新交易，在以太坊虚拟机中执行它们，并保存所有当前以太坊数据的最新状态和数据库。 

    由**Geth**(Go) 客户端主导，OpenEthereum (Rust)客户端远远排在第二，Erigon (Go) 客户端排第三，Nethermind (C#、.NET Core) 排第四，其他客户端只占不到网络的 1%，如Besu (Java)等。

  - **共识客户端**：也称为信标节点，实现了权益证明共识算法，使网络能够根据来自执行客户端的经过验证的数据达成一致。

    最常用的客户端 **Prysm** (Go)虽然不像执行层的 Geth 客户端那样占主导地位，但仍然拥有了超过 60% 的网络，Lighthouse (Rust)和 Teku (Java)分别占 20% 和 14%，其他客户端则很少使用。

- *以太坊节点*

  - 归档节点：
    - 保存区块链的所有内容，并建立历史状态存档。 如果你想查询任意区块高度的帐户余额，或者想本地测试自己的一组交易而不必真的mined，则需要使用归档节点通过trace获取。
    - 这些数据TB为单位，这使得归档节点对普通用户的吸引力较低，但对于区块浏览器、钱包供应商和链分析等服务来说却很方便。

  - 全节点：
    - 存储全部区块链数据（会定期修剪，源码default只保留最近128个区块全状态，并不存储包含创世块在内的所有状态数据）
    - 参与区块验证，验证所有区块和状态。
    - 所有状态都可以从全节点中获取（未存储的状态是通过向归档节点发出请求重建的）。
    - 为网络提供服务，并应要求提供数据（比如帮助轻节点验证）

  - 轻节点：

    - 轻节点仅下载区块头， 所需的任何其他信息都从全节点请求，然后可以根据区块头中的状态根独自验证收到的数据。这帮助用户无需装备昂贵的硬件或高带宽，就可以加入以太坊网络。

    - 轻节点不参与共识（即不能成为矿工或验证者），但可以访问以太坊区块链，其功能与全节点相同。
    > 以太坊有一个LES协议，light ethereum subprotocol，是为轻节点通讯准备的。执行这个协议还是靠自愿的，没有强制的措施让客户端去执行这个协议，如果你运行了一个全节点，并且愿意服务轻节点的话，比如geth，可以打开 light.serve，限制自己执行LES的时间比例，也可以用light.ingress, egress, maxpeer等命令，限制自己为轻节点付出的负担。多数做轻节点客户端的团队，自己也会跑全节点来运维与测试。

- *以太坊中的团体与组织*

  - [**以太坊基金会**](http://ethereum.foundation/) (Ethereum Foundation, EF) 是一个非营利性组织，致力于支持以太坊以及相关技术。EF 不是一家公司，甚至不是传统的非营利组织。 他们的作用不是控制或领导以太坊，也不是为与以太坊相关的关键技术开发提供资金的唯一组织。 EF 只是巨大的生态系统的一部分。

  - [**以太坊核心开发者**](https://mp.weixin.qq.com/s/cKqT18yRu4dBsKULFE7z5Q)(Ethereum Core Developer)，指的是正在（Currently）为以太坊底层协议开发提供重要贡献的人。

    > 比如，虽然以太坊联合创始人 Gavin Wood 曾经为早期的以太坊作出重大贡献，他现在已经不再被认为是以太坊核心开发者了，只是前核心开发者。

    **包括各个客户端的开发者，底层协议的开发者及核心的以太坊研究员**(同时也是这些人受邀参加核心开发者会议)。各客户端的开发者是核心开发者的子集，是以太坊众多客户端开发者中的一部分。以太坊会有自己的技术规范，各个客户端的开发者对根据这个规范来开发客户端。

#### 正文）

无论是通过 `geth()` 函数还是其他的命令行参数启动节点，节点的启动流程大致都是相同的，这里以 `geth()` 为例`/cmd/geth/main.go`

=> Init() 使用了 `gopkg.in/urfave/cli.v1` 扩展包，通过初始化了命令行解析配置

=> main()启动app程序并解析

=> geth() 配置节点相关服务backend，启动以太坊节点node

> `Node` （`node/node.go`）是 geth 生命周期中最顶级的实例，负责作为与外部通信的高级抽象模块的管理员，比如管理 rpc server，http server，Web Socket，以及P2P Server外部接口。同时，Node中维护了节点运行所需要的backend的实例和服务(`lifecycles []Lifecycle`)。
>
> `Backend `（`eth/backend.go`）中`Ethereum`结构体包含的成员变量以及接收的方法实现了一个Ethereum full node所需要的全部服务功能和数据结构。我们可以在下面的代码定义中看到，Ethereum结构体中包含了`TxPool`，`Blockchain`，`consensus.Engine`，`miner`等最核心的几个数据结构作为成员变量。

```go
func geth(ctx *cli.Context) error {
   if args := ctx.Args(); len(args) > 0 {
      return fmt.Errorf("invalid command: %q", args[0])
   }

   // 1.选择对应以太坊网络，预分配内存缓存，启动服务监控
   prepare(ctx)
  
   // 2.stack为Node实例，backend为节点运行所需的后端实例，提供更为具体、底层的以太坊的功能性Service
   // makeFullNode => makeConfigNode返回default节点，后续导入配置至该节点
   stack, backend := makeFullNode(ctx)
   defer stack.Close()
  	
   // 3.启动一个以太坊节点Node.Start()
   // 会从Node.lifecycles中注册的backend服务实例，并启动它们
   startNode(ctx, stack, backend, false)
   
   // 堵塞主线程，其他的功能模块的服务被分散到其他的子协程中进行维护,通过内置函数close()关闭
   stack.Wait()
   return nil
}
```

![img](https://images.seebug.org/content/images/2018/07/2d428e8a-9276-4645-88bc-20f14a125d01.png-w331s)

