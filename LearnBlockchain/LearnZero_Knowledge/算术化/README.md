# 算术化

**算术化**是将计算编码为代数约束满足问题的过程。这将检验其正确性的复杂性降低到少量概率代数检查。在证明系统中，算术化的选择会影响IOP的选择范围

<div align=center>
<img src="https://zkshanghai.xyz/assets/components_of_proof_system_zh.drawio.854a14bc.svg" style="width:20%;">
</div>
## 一、高效密码学运算

### a.离散傅里叶变换

[傅里叶变换](https://www.jezzamon.com/fourier/)可以将一个信号分 解成一系列正弦和余弦函数的叠加。

离散傅里叶变换(DFT)：

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/3bab2249-2b01-4b32-af11-9c449fd0fb08" style="width:60%;">
</div>


其中一个应用场景：每个约束都是一个多项式，可利用DFT做到整合成一个多项式且保证信息不丢失。

### b. 单位根(Root of unity)

- 称 $𝑥^𝑛 = 1$ 在复数意义下的解是 𝑛 次复根。
- 这样的解有 𝑛 个，称这 𝑛 个解都是 𝑛 次 单位根 或 单位复根 （the 𝑛-th root of unity）。 
- 根据复平面的知识，𝑛 次单位根把单位圆 𝑛 等分。 

> 如复平面所示为 $x^3 -1$ 的三次单位根等分情况（ $x=1$ 旋转 $0$, $\frac{2π}{3}$, $\frac{4π}{3}$ ）

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/e4da36bf-7392-4019-aa47-05dd739535ba" style="width:20%;">
</div>

### c.典型数学难问题类型

> **数学难问题**知识补充：
>
> **P**： 表示可以由确定性图灵机，**确定在多项式时间内解决**的判定问题，这也是目前经典计算机的运算能力。
>
> [**NP**](https://zhuanlan.zhihu.com/p/73953567)：由[非确定性图灵机](https://www.cnblogs.com/zhangzefei/p/9742918.html)可以在多项式时间内解决的判定问题，**不确定**有没有多项式解决算法，但可以通过验证的方法得出正确解，所有的P问题都是NP问题。
>
> **NP-Hard**：如果所有 NP 类问题都可以在多项式时间内规约到问题 H，那么问题 H 是 NP-hard 的。**NP-Complete**：如果一个问题，既是NP类问题，又是NP-hard问题，这是**零知识证明的基础**。

> 难度：DLP > CDH > DDH

#### 1.[离散对数问题（The Discrete Logarithm Problem, **DLP**）](https://chenliang.org/2021/05/09/discrete-logarithm-problem-and-DHKE/)

> 离散：有限域而非实数域，但需要注意，不是所有的离散对数问题都是困难的。
>

应用：

- 椭圆曲线


<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/e7b54d01-40b8-42f4-b0f7-216de220cddf" style="width:60%;">
</div>



- [DH密钥交换协议](https://chenliang.org/2021/05/09/discrete-logarithm-problem-and-DHKE/)，在不安全的通道，通过 shared secret 建立安全的传输。

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/50f715a5-70d8-4ae3-84d2-59fb3fc69af4" style="width:60%;">
</div>

- [Schnorr Protocol](https://crypto.stackexchange.com/questions/9997/perfect-zero-knowledge-for-the-schnorr-protocol)

  符合完美零知识(见下文区别中介绍)，最后的 $z$ 中， $r$ 的随机性保护了 $s$ 。Schnorr的签名方案是一个经典的Sigma协议，具有Special Honest-Verifier Zero-knowledge property。

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/6a4a961e-edef-4eed-85a0-e481cc586fad" style="width:60%;">
</div>

#### 2.计算性DH问题（The Computational Diffie-Hellman Problem, **CDH**）

属于功能性问题，功能性问题的回答不止 YES/NO，可以是一个数或是其它。如「求两个数的和」就是一个功能性问题。

#### 3.判定性DH问题（The Decisional Diffie-Hellman Problem, **DDH**）

只能用 YES/NO 回答的问题，本质上是判断是否属于某一种**语言**(是一个抽象的概念，通常意义上的字符串是语言，所有的有向无环图也可以是一个语言)。



## 二、例1：R1CS到QAP

![image-20240122204153366](/Users/dazso/Library/Application Support/typora-user-images/image-20240122204153366.png)

### a.一阶约束系统 (Rank-1 Constraint System, R1CS)

本质是一个方程组，由一个 R1CS 是一个由三个向量构成的向量组 $(\vec{a},\vec{b},\vec{c})$ ，假设 R1CS 的解也是一个向量，记为 $\vec{x}$ ，其中 $·$ 表示向量内积， $∗$ 表示算数乘法：
$$
(\vec{x}·\vec{a}) * (\vec{x}·\vec{b}) =(\vec{x}·\vec{c})
$$

### b.二次算术程序（Quadratic Arithmetic Program，QAP)

是一种**将语句转换为多项式上二次方程组的方式**，它们可以通过线性交互式证明LIPs、代数IOPs、多线性IOPs 等不同信息论协议进行检验。

QAP 实现了与 R1CS 完全相同的逻辑，只不过使用的是多项式而不是向量内积。任何具有乘性复杂度 $n$ 的电路都可以转换为一个 $n$ 次多项式的QAP。通过QAP 转换，可以将算术电路（计算问题）转换为二次多项式形式，等式中的每个变量至多为二次，形式为：
 $$\sum_{k}{A_{𝑖k}z^k} * \sum_{k}{B_{𝑖k}z^k} = \sum_{k}{C_{𝑖k}z^k}$$ 
**QAP判定**：

一个度数（多项式最高次）为 $d$ 、大小为 $m$ 的二次算术程序 $Q$ 由多项式， 和一个目标多项式 $T(X):=\prod_{i = 0}^{d-1}{(x-i)}$  组成。当赋值 $(1,x_1,…,x_{m−1})$ 满足 $Q$ 时，有 $$𝑃(𝑋）:= 𝐿(𝑋) · 𝑅(𝑋) − 𝑂(𝑋)$$，其中

> 丨：整除

$$
𝑇(𝑋) ∣ 𝑃(𝑋)
$$

此外，在转换的同时会构建一个对应于**代码输入的解（又称为 QAP 的 Wit­ness）**，之后再基于这个 Wit­ness 构建一个实际的零知识证明系统

### c.实现步骤

> 以IsZero()判零函数（见算术点路应用章节a.①）为例，或从`comparators.circom`中获取的`IsZero`电路。

#### 1.Flattening扁平化（通过代码构建树，将算术点路扁平化）

`out <== −in * inv + 1`

`in * out === 0`


<div align=center>
<img src="https://zkshanghai.xyz/assets/iszero_circuit_dag.9e49a073.png" style="width:50%;">
</div>

“展平” 为四个约束条件(在算术电路表示中，每个约束条件对应于一个加法或乘法门)，每个都采用 `左侧 · 右侧 = output` 的形式：

> g：门gate

- $𝑔_0： 𝑤_1 · −1 = 𝑤_2$ 
- $𝑔_1：𝑤_2 · 𝑤_3 = 𝑤_4$ 
- $𝑔_2： (𝑤_4 + 1) · 1 = 𝑤_5$ 
- $𝑔_3： 𝑤_1 ⋅ 𝑤_5 = 𝑤_6$

#### 2.结果转换为**R1CS**（从计算问题转为矩阵式）

为满足 $L\vec{x}·R\vec{x}=O\vec{x}$ 形式，创建三个线路向量 $\vec{l_i}$、 $\vec{r_i}$ 、 $\vec{o_i}$  ，其中包含门中每个变量 $\vec{w_i}$ 的系数，线路向量还包括一个常数项 $\vec{w_0}$ ，即算上常数项共七列 $(w_0,w_1, w_2,w_3, w_4,w_5, w_6)$

<img width="396" alt="image" src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/3d23dc54-7476-45f9-a60a-cfb500c9d1e6">


分别收集三个线路向量到矩阵 $L$， $R$ ， $O$ ，与见证向量(Witness，证明人声称知道一些合法的赋)  $\vec{x} = (1,x_1, x_2, x_3, x_4, x_5, x_6)$  一起，构成了判零电路的R1CS形式。

$$
L=
\left(
\begin{matrix}
\vec{l_0}:0 & 1 & 0 & 0 & 0 & 0 & 0\\
\vec{l_1}:0 & 0 & 1 & 0 & 0 & 0 & 0 \\
\vec{l_2}:1 & 0 & 0 & 0 & 1 & 0 & 0\\
\vec{l_3}:0 & 1 & 0 & 0 & 0 & 0 & 0
\end{matrix}
\right)
$$

$$
R=
\left(
\begin{matrix}
\vec{r_0}:-1 & 0 & 0 & 0 & 0 & 0 & 0\\
\vec{r_1}:0 & 0 & 0 & 1 & 0 & 0 & 0 \\
\vec{r_2}:1 & 0 & 0 & 0 & 0 & 0 & 0\\
\vec{r_3}:0 & 0 & 0 & 0 & 0 & 1 & 0
\end{matrix}
\right)
$$

$$
O=
\left(
\begin{matrix}
\vec{o_0}:0 & 0 & 1 & 0 & 0 & 0 & 0\\
\vec{o_1}:0 & 0 & 0 & 0 & 1 & 0 & 0 \\
\vec{o_2}:0 & 0 & 0 & 0 & 0 & 1 & 0\\
\vec{o_3}:0 & 0 & 0 & 0 & 0 & 0 & 1
\end{matrix}
\right)
$$

#### 3.**R1CS转为QAP**（矩阵式转为多项式）

> 不管结果多简单，过程是复杂的，要验算的东西即组成的多项式也很复杂，所有的中间变量都会成为输入。
>
> 注：circom是从代码到R1CS，他下面还有个 snarkjs工具框架，将R1CS生成QAP

将度数 $d$ 视为约束的数量，将规模 $m$ 视为变量的数量。在我们的例子中，有 $d=4$, $m=7$ 。通过将R1CS形式转换为QAP形式，从三次矩阵乘法降低到单项式恒等式。

在每个变量 $j$ 和门 $i$ 处，我们希望 $L_j(i)$ 选择门 $g_i$ 的左导线的变量 $w_j$ 的系数；$R_j(i)$ 和 $O_j(i)$ 同理，根据QAP判定式判断是否成立。

$L_j(i) = L_{ij} = \vec{l}_i[j]$

$R_j(i) = R_{ij} = \vec{r}_i[j]$

$O_j(i) = O_{ij} = \vec{o}_i[j]$

其中，在构造 $L_j$ 时，我们将每个 $L_j$ 设置为在评估点 $(0,…,d−1)$ 列 $L[j]$ 中的值的插值多项式；  $R_j$， $O_j$  同理。



**拉格朗日插值**：

> 注意实际运行时有限域，无小数

1. 找到一个多项式，经过四个关键点
2. 当处在某个控制点 $x_i$ 的情况下，除了该点有值，其他的控制点值为0

实际是多个多项式相加而得，$y_i$ 为缩放系数，给定点和评估 ${(x_i,y_i)}_{i=0}^{d−1}$ ，我们可以构造一个**插值多项式** $L(X)$，使 $L(x_i)=y_i$ :

 $$L(X):=\sum_{i=0}^{d-1}{y_i·L_i(X)}$$ 
其中， $L_i(X)$ 是穿过评估点{ ${x_0,…,x_{d−1}}$ } 的拉格朗日基本多项式：

<img width="369" alt="image" src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/a3a61b14-65bf-4900-9db5-5d8e9d430256">



*注：得出一个多项式，就成了上节承诺方案的输入*



## 三、例2：代数中间表示AIR

> R1CS生成不一定是QAP，根据不同框架，这里是到AIR，相当于分别是不同的算术化方法

**代数中间表示（Algebraic Intermediate Representation, AIR）**，是由一组**均匀计算**（uniform computations，均匀性是指数据或分布的相似程度。在统计学中，均匀性的计算通常是指方差或标准差）组成的程序表示，是 StarkWare 在其虚拟机 Cairo 中使用的算术化过程。

### a.特点

> AIR框架数学上实现简单，构建出那么一个表格（轨迹矩阵，像机器语言底层系统，故也有的称AIR为机器计算），所有多项式乘在一起就得出唯一多项式。

1. 计算**执行轨迹**。表示为执行迹矩阵 T ，**行**表示在给定时间点的计算状态，**列**对应于一个代数寄存器在所有计算步骤中的状态变化。
2. **转换程序**。约束了迹矩阵 T 两行或多行之间的关系。
3. **边界约束**。确保了执行中某些单元格和特定常量之间的相等关系。

### b.以Fibonacci为例论证过程

#### 1.执行轨迹Table

> 表格高度是固定的，像circom在运行时不可变更（树的高度固定）

| Step |  a   |  b   |
| :--: | :--: | :--: |
| i=1  |  1   |  1   |
| i=2  |  2   |  3   |
| i=3  |  5   |  8   |
| i=4  |  13  |  21  |

#### 2.转换程序

我们可以使用两个状态转换多项式来指定 Fibonacci 数列的 AIR 程序：

$f_1(X_1,X_2,X_1^{next} ,X_2^{next})=A^{next}-(B+A)$

$f_1(X_1,X_2,X_1^{next} ,X_2^{next})=B^{next}-(B+A^{next})$

#### 3.边界约束（结果应 = 0）

例如，我们可以检查在第 i=2 行状态转换是否成立：

$f_1(X_1,X_2,X_1^{next} ,X_2^{next})=5−(3+2)=0$

$f_2(X_1,X_2,X_1^{next} ,X_2^{next})=8−(5+3)=0$

### c.PAIR, Preprocessed Algebraic Intermediate Representation

带预处理的AIR，以在转换程序中同时启用乘法和加法。

### d. RAP, Randomized AIR with Preprocessing

**带预处理的随机化AIR**，以实现**多重集合相等性检查**。允许交互轮次引入验证器随机性 $γ$ ∈ 有限域 $F$。在稍后的轮次中，可以将较早轮次的随机性用作约束中的变量。这使得本地约束（相邻行之间）可以检查全局属性。

### e.利用AIR构造ZKVM

> 针对AIR框架，STARK、SNARK（一开始NIZK用得最多，后面加的succint性质）是人为定义的schema，是一种性质，不是很具体的框架，STARK提出者(公司STARKNET)底层用的AIR，同时可以认为STARK也属于SNARK一种实现方式
>
> 通常见到AIR就是已经封进了ZKVM中，可以实现传统代码，像StarkNet的Cairo，开放给用户来写代码，不涉及到AIR底层

**算术化步骤：**

将待验证计算问题转换为检查某个多项式，分两步：

-  第一步

  -  构建**执行踪迹表格**

    <img src="/Users/dazso/Library/Application Support/typora-user-images/image-20240122210357730.png" alt="image-20240122210357730" style="zoom: 33%;" />

  -  用多项式描述表格中各行/列间的数学关系 

- 第二步： 将这两个对象转换为一个低次多项式

  - 利用**纠错码**将执行轨迹转为多项式

    > 哪怕仅一处错误的执行轨迹，会被纠错码放大，以至于与原执行轨迹几乎完全不同

    <img src="/Users/dazso/Library/Application Support/typora-user-images/image-20240122210417088.png" alt="image-20240122210417088" style="zoom: 50%;" />

     

  - 扩展至更大的域 

  - 用多项式约束将其转为低次多项式

