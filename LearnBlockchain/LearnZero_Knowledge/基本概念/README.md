![image](https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/37c21de1-54e0-4b8f-a19c-877d30c55857)# 零知识证明基本概念

## 一、定义

下图为 [Mental Poker over the Telephone[SRA81]](https://people.csail.mit.edu/rivest/pubs/SRA81.pdf)，文中提出了通过完备的论证方法，来实现对传递信息的隐匿，而非严格的数学证明，该结论为后续零知识证明的可行性奠定了基础。2012年图灵奖获得者Shafi Goldwasser与Silvio Micali首次正式引入了零知识证明（Zero-Knowledge Proof, ZKP）的概念，文章为[The knowledge complexity of interactive proof-systems](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Proof%20Systems/The_Knowledge_Complexity_Of_Interactive_Proof_Systems.pdf)

> **零知识证明，实际都是通过论证得出，而非严格的数学证明**
>
> 论证：利用论据（原因，例子，数据）等，验证论点为真（无限接近100%）。交互证明(Interactive Proof)在定义上是允许一个无效证明通过验证的可能性微小但非零概率，而论证即为计算可靠性问题的IP，允许存在不正确陈述的“证明”（需要非常大的算力才能找到）。以计算性问题为前提条件后，一个重要的优势是能够使用(多项式时间安全)的密码学原语。
> 
> 证明：从空条件集合出发，利用假设法、公理、公设，来证明一个命题的恒为真，每个过程都要详细，确信的100%

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/0c0401f8-31b6-4310-b7eb-97dfae251532" style="width:65%;">
</div>

**定义** :零知识证明（ZKP, Zero-Knowledge Proof）是**一方（证明者Prover）向另一方（检验者Verifier）论证某命题的方法**。零知识协议是一种协议，它允许你证明你知道一些特定的数学事实，而无需透露有关事实本身的任何信息。在零知识协议中生成的证明称为零知识证明。

- 完备性(completeness)：如果 Prover 知道正确的具体信息，他们总是能够做出令人满意的回答；

- 可靠性(Soundness)：如果 Prover 不知道正确的具体信息，他们最终会被抓住，换言之，可靠的系统不会允许虚假的陈述被错误地证明为真;

  > 知识可靠性(Knowledge Soundness)：是可靠性的更强要求，**要求证明者必须真的知道某个特定的信息（比如一个密码或是一个解）**。（该点往往是使用者实现，比如有的ZK-Snark系统本身只能是实现可靠性，因为系统内不包含“知识”信息，开发者在该系统上架构，可以实现知识可靠性）

- 零知识(Zero Knowledge)：Prover 的回答不会泄露其具体信息，验证者除了正确外，不知道任何其他信息；

  > 值得说明的是，零知识特性是**可选的(optional)**。例如，zkrollup的实现过程中，一般只会实现SNARK(见ZK-SNARKs章节)，而不会引入ZK属性


**应用**：零知识证明都是针对**可计算问题**

> 对于一个判定问题，若存在一个总是在有限步内停机且能够正确进行判定的图灵机，则这个问题是一个 **图灵可计算** 的问题，否则这个问题是一个 **图灵不可计算** 的问题。

> 典型不可计算问题-停机问题：给定 $α$ 和 $x$ ，判定 $M_α$ 在输入为 $x$ 时是否会在有限步内停机。

## 二、现代零知识的特点

1. 为解决通用可计算问题，设计一个**通用框架将问题转化为一个可计算模型(computational model)**，然后利用零知识证明系统生成证明。在古典零知识证明中，对于像[三色图问题](https://www.jianshu.com/p/7b772e5cdaef)，研究人员（数学家）都必须提出一个应用于特定场合的特殊用途/特定的ZK协议，而在现代ZK零知识证明中，。

   > 一个输出 $y$ ，一个任意函数 $f$，知道一个秘密 $x$ ，使得 $f(x)=y$

   在此基础上，我们能够在完全隐私的情况下验证任意计算（如货币、数字所有权的转移）。也是因此，程序员才可以上手写逻辑。

2. **满足完美零知识证明的不可区分特性**。如果 “[模拟](https://sammyne.github.io/zkp/02-simulation/#%E5%8C%BA%E5%88%86%E4%B8%A4%E4%B8%AA%E4%B8%96%E7%95%8C)视图” 和 “真实交互” 通过**任意算法**无法区分，则称为完美零知识。

   > 计算不可区分加密：限制了任意算法为 多项式时间算法(ploy-time algorithm)

   例如，在现代零知识证明中，非对称加密(如椭圆曲线)的数字签名具**不具备**零知识属性。

   > 如，私钥在一个向量中[s1,s2,s3,...,sN]，使用一个普通的椭圆曲线签名可以证明你拥有某个私钥sK，但是一个adversary可以通过遍历这个向量来模拟出你的信息，发现你用的是哪个私钥，零知识需要保证哪怕这个adversary 能够模拟出一些信息,也无法得知用哪个。
   >
   > 进阶：Schnorr签名方案，一个经典的**Sigma协议**(见承诺方案章节)，具有Special Honest-Verifier Zero-knowledge property

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/c01c0228-f0d7-4dbd-be7d-90b7c8c098f6" style="width:65%;">
</div>

