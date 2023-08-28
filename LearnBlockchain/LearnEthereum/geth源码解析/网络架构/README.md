### [网络架构](https://paper.seebug.org/642/#0x02-geth)

#### 1）逻辑架构

在以太坊中，p2p 作为通信链路，用于负载上层协议的传输，可以将其分为三层结构：

![img](https://images.seebug.org/content/images/2018/07/0f5a01a4-63e9-4d72-bbf1-b492d52379cd.png-w331s)

1. 最上层是**以太坊中各个协议的具体实现**，如 eth 协议、les 协议。
2. 中间层即 **p2p 通信链路层**，主要负责启动监听、处理新加入连接或维护连接，为上层协议提供了信道。
3. 最下面的一层，是由 Go 语言所提供的**网络 IO 层**，也就是对 `TCP/IP` 中的网络层及以下的封装。

**其中，P2P通讯链路层主要司职以下工作：**

![img](https://images.seebug.org/content/images/2018/07/18197a44-c04e-4650-b843-3d87ce4f5700.png-w331s)

1. 由上层协议的数据交付给 p2p 层后，首先通过 **RLP 编码**(可见上文编码方式)。
2. RLP 编码后的数据将由**共享密钥加密**（`/p2p/transport/doEncHandshake()`），保证通信过程中数据的安全。

> **「迪菲-赫尔曼密钥交换」**
> p2p 网络中使用到的是（英语：Diffie–Hellman key exchange，缩写为D-H） 是一种安全协议。它可以让双方在完全没有对方任何预先信息的条件下通过不安全信道创建起一个密钥。
>
> 简单来说，链接的两方生成随机的私钥，通过随机的私钥得到公钥。然后双方交换各自的公钥，这样双方都可以通过自己随机的私钥和对方的公钥来生成一个同样的共享密钥(shared-secret)。后续的通讯使用这个共享密钥作为对称加密算法的密钥。`ECDH(A私钥, B公钥) == ECDH(B私钥, A公钥)`。

3. 将数据流**转换为 RLPXFrameRW 帧**（`/p2p/transport/readMsg()&Write()`），便于数据的加密传输和解析。

> 目的：在单个连接上支持多路复用（Multiplexing）协议，以便高效传输
>
> 帧结构：五个数据包
>
> ```go
> header          // 包含数据包大小和数据包源协议
> header_mac      // 头部消息认证
> frame           // 具体传输的内容
> padding         // 使帧按字节对齐
> frame_mac       // 用于消息认证
> ```

#### 2）[ÐΞVp2p源码分析](https://learnblockchain.cn/article/1937)

> 以太坊所实现的P2P协议成为：devp2p

（`/cmd/p2p/server.go`）

启动服务目录及堆栈：

> /cmd/geth/main.go => geth() => startNode() => utils.StartNode() => stack.Start() => openEndpoints() => **server.Start()**

**[启动p2p网络细节](https://learnblockchain.cn/article/1937)，主要会做以下几件事：**

> 1. 初始化server的字段
> 2. 设置本地节点setupLocalNode
> 3. 设置监听TCP连接请求setupListening
> 4. 设置节点发现（setupDiscovery）, 利用[KAD算法](https://zhuanlan.zhihu.com/p/43340851)，其[源码实现分析](https://blog.csdn.net/lj900911/article/details/84138361)
> 5. 设置最大可以主动发起的连接遵循50/3规则
>    > 最大主动连接数量：这意味着节点会尝试主动连接到其他节点，直到它达到这个限制。在许多以太坊客户端中，这个值默认是 50。
>    > 
>    > 保留（reserved）：这个值表示即使达到了最大主动连接数量，还有多少连接是“保留”的，可以被特定的节点（例如，你事先知道并信任的节点）使用。在很多以太坊客户端中，这个值默认是 3。
> 7. srv.run(dialer) 发起建立TCP连接请求

```go
func (srv *Server) Start() (err error) {
   
   // 读写锁与p2p服务参数初始化
 	 ......

   // 1.启动一个以太坊节点记录ENR(EIP778)形式本地节点 
   // 以太坊网络地址的标准格式, 取代了 multiaddr和enode, 使节点之间能够进行更多的信息交流
   if err := srv.setupLocalNode(); err != nil {
      return err
   }
  
   // 2.服务监听
   if srv.ListenAddr != "" {
      if err := srv.setupListening(); err != nil {
         return err
      }
   }
  
   // 3.根据用户参数开启`节点发现`功能，基于kademlia（KAD）算法
   // 使本地节点得知其他节点的信息，进而加入到p2p网络中
   // 注意：KAD算法需要种子节点来引导，从leveldb中随机选取若干种子节点（新节点第一次启动时，使用启动参数或源码中提供的启动节点作为种子节点）
   // 其实现细节位于：p2p/discover/table.go
   if err := srv.setupDiscovery(); err != nil {
      return err
   }
  
   // 4.初始化dialer，负责和peer建立连接关系
   srv.setupDialScheduler()
   srv.loopWG.Add(1)
  
   // 5.单独协程进行报文处理
   go srv.run()
   return nil
}
```

![img](https://images.seebug.org/content/images/2018/07/294c3853-6aa4-48a8-8f95-8fec8a17b568.png-w331s)
