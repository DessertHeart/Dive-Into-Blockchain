# 🧮Learn Zero-Knowledge

ZK是密码学的一个分支，是一种证明方法，也可以称为技术。其概念的热度，自2022开始持续至今，无论事今年3月份结束的 ETH Denver 大会上还是4月份香港万向Web3 Festival中，**ZK 俨然成为开发者和投资者们交流中最高频的热词**。2023年 5 月 19 日- 23 日黑山举办的 EDCON 2023 大会上，以太坊联合创始人 Vitalik Buterin（V神）更是直言，未来 10 年 **ZK-SNARK 将与区块链一样重要** 。由此可见，学习和掌握ZK技术，对于Blockchain Developer的未来发展有着不可忽视的作用。

本文档基于 [Sutulabs 和 Kepler42B - ZK Planet 共同举办的零知识开发者Workshop](https://zkshanghai.xyz/)与[[MIT IAP 2023] Modern Zero Knowledge Cryptography](https://zkiap.com/)进行课堂笔记的整理与知识点汇总，旨在帮助各位区块链开发者更好的学习ZK技术，并提供平台与大家一起交流探讨！

> [课后作业](https://github.com/DessertHeart/zkshanghai-workshop/tree/main)

+ [零知识证明基本概念](./%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
+ [承诺方案Commitment](./%E6%89%BF%E8%AF%BA%E6%96%B9%E6%A1%88)
  1. Pedersen承诺 
  2. 向量Pedersen承诺
  3. 双线性映射
  4. KZG承诺
+ [算术化](./%E7%AE%97%E6%9C%AF%E5%8C%96)
  1. 从R1CS到QAP
  2. AIR
  3. 高效密码学运算
+ [ZK-SNARKs](./ZK-SNARK)
  1. 特点
  2. 实现过程
  3. PLONK协议
+ [算术电路应用与Circom](./%E7%AE%97%E6%9C%AF%E7%94%B5%E8%B7%AF%E5%BA%94%E7%94%A8%E4%B8%8ECircom)
  1. 基础电路
  2. 复杂实用电路
  3. Circom
  4. 应用ZK结构
  5. 应用ZK实例分析



> 附录：👨‍🎓由浅入深学习零知识证明资料推荐

- [安比实验室零知识证明小白入门](https://www.zhihu.com/question/265112868/answer/891098212)；
- [安比实验室的工作，介绍从0到1实现Pinocchio算法](https://zhuanlan.zhihu.com/p/99260386)，著名的Groth16算法的前身，从Pinocchio算法了解如何实现一个ZKP系统；
- [ZK-Learning Mooc 课程](https://zk-learning.org/)，大牛视频授课(Youtube & B站)，覆盖了需要了解的数学知识以及ZKP系统构造；
- [UCBerkeley的菲尔兹奖得主 Richard Borcherds 录的数论课程](https://www.youtube.com/playlist?app=desktop&list=PL8yHsr3EFj53L8sMbzIhhXSAOpuZ1Fov8)，深入学习零知识证明所需数学知识；
- ["Proof, Argument and Zero-knowledge" ](https://people.cs.georgetown.edu/jthaler/ProofsArgsAndZK.pdf)书籍，学术硬核，深入学习零知识证明算法实现；
- 研读零知识证明的各种构造论文……
