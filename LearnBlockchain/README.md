# :books: LearnBlockchain

深入公链技术、跨链技术等底层技术应用的学习手册，从技术开发角度掌握其原理

## LearnEthereum -精通以太坊

+ [Dapp开发入门](./LearnBlockchain/LaernEthereum/Dapp%E5%BC%80%E5%8F%91%E5%85%A5%E9%97%A8)
+ [Ethereum架构](./LaernEthereum/EVM%E5%BA%95%E5%B1%82%E7%BB%93%E6%9E%84)
+ go-ethereum源码剖析
  + 1、[geth目录结构](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/geth%E7%9B%AE%E5%BD%95%E7%BB%93%E6%9E%84)
  + 2、[以太坊初始化](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E5%88%9D%E5%A7%8B%E5%8C%96)
  + 3、[P2P网络架构](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E7%BD%91%E7%BB%9C%E6%9E%B6%E6%9E%84)
  + 4、[Account账号结构](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E8%B4%A6%E6%88%B7%E7%BB%93%E6%9E%84)
  + 5、[世界状态与stateDB](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E4%B8%96%E7%95%8C%E7%8A%B6%E6%80%81State%E4%B8%8EStateDB)
  + 6、[交易Transaction](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E4%BA%A4%E6%98%93Transaction)
  + 7、[共识协议Block => Blockchain](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/%E4%BB%8EBlock%E5%88%B0Blockchain)
  + 8、[Log与布隆过滤器](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/Log%E5%92%8C%E5%B8%83%E9%9A%86%E8%BF%87%E6%BB%A4%E5%99%A8)
  + 9、[EVM与Opcodes](./LaernEthereum/geth%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/VM%E5%92%8COpcodes)

## Learn Zero-Knowledge -入门零知识证明

本文档基于 [Sutulabs 和 Kepler42B - ZK Planet 共同举办的零知识开发者Workshop](https://zkshanghai.xyz/)与[[MIT IAP 2023] Modern Zero Knowledge Cryptography](https://zkiap.com/)进行课堂笔记的整理与知识点汇总，旨在帮助各位区块链开发者更好的学习ZK技术，并提供平台与大家一起交流探讨！


+ [零知识证明基本概念](./LearnZero_Knowledge/%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
+ [承诺方案Commitment](./LearnZero_Knowledge/%E6%89%BF%E8%AF%BA%E6%96%B9%E6%A1%88)
  1. Pedersen承诺 
  2. 向量Pedersen承诺
  3. 双线性映射
  4. KZG承诺
+ [算术化](./LearnZero_Knowledge/%E7%AE%97%E6%9C%AF%E5%8C%96)
  1. 从R1CS到QAP
  2. AIR
  3. 高效密码学运算
+ [ZK-SNARKs](./LearnZero_Knowledg/ZK-SNARK)
  1. 特点
  2. 实现过程
  3. PLONK协议
+ [算术电路应用与Circom](./LearnZero_Knowledg/%E7%AE%97%E6%9C%AF%E7%94%B5%E8%B7%AF%E5%BA%94%E7%94%A8%E4%B8%8ECircom)
  1. 基础电路
  2. 复杂实用电路
  3. Circom
  4. 应用ZK结构
  5. 应用ZK实例分析
