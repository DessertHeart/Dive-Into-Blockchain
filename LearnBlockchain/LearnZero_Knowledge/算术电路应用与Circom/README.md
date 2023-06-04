# 算术电路应用与Circom

## a.基础电路

**定义** :在计算复杂性理论中，计算多项式的一个计算模型。对于一个给定的有限域 $F$ = { ${0, … , 𝑝 − 1}$ } ，基于某一素数 $𝑝 > 2$，算术电路 $𝐶$:  计算一个在 $F[x_1,...,x_n]$ 中的多项式。

> 多项式时间复杂度：指解决问题需要的时间与问题的规模之间是多项式关系。形似  $O(nlog(n))、O(n^2)、O(n^3)$

- 算术电路是计算复杂性理论中的概念，与电子电路毫无关联
- 有向无环图
- 输入节点标记为 $1, 𝑥_1, … , 𝑥_𝑛$
- 内部节点标记为 $+,-,x$
- 每个**内部节点也称为门 (gate)**

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/fcfe3568-f84a-4876-96cb-abccc054cc91" style="width:25%;">
</div>

**特点：**

1. 固定性：电路在证明过程不可动态增减
2. 丰富性：电路传递的是有限域的数字，比二进制具有丰富的表达能力
3. 约束性：电路既可用作计算，也用于约束流转信号的状态转换

**零知识证明电路常见设计方法**

> **设计的原因**：因为目前零知识领域正在使用 “算术电路” 这样的构造方式，与我们的传统程序逻辑不一样，因 此过去那些成熟的程序逻辑，没办法轻易的移植到零知识证明器里，故很多逻辑**由程序员根据实际场景进行电路设计**。

> **提示$hint$**：加速计算。出于算术电路的约束性，每个门电路的计算都会转换为约束，进而增加证明和验证的工作量，我们可以将复杂的计算过程变为预先计算的提示值，在电路中对提示值进行验证，从而降低证明和验证的工作量。用于验证一些类似 $A$ != 0。

### ①判零函数

`let y = input > 0 ? 0 : 1;`

> 如果用[费马小定理 (判断一个数是不是素数)](https://zhuanlan.zhihu.com/p/87611586)计算量非常大。

利用提示值 𝑖𝑛𝑣 计算输出𝑜𝑢𝑡𝑝𝑢𝑡 并验证输入输出符合约束。𝑖𝑛𝑝𝑢𝑡 = 0，提示值 𝑖𝑛𝑣 =  0；否则 𝑖𝑛𝑣 = 1/𝑖𝑛𝑝𝑢𝑡

1. $𝑜𝑢𝑡𝑝𝑢𝑡 = −𝑖𝑛𝑝𝑢𝑡 × 𝑖𝑛𝑣 + 1$
2. $𝑖𝑛𝑝𝑢𝑡 × 𝑜𝑢𝑡𝑝𝑢𝑡 = 0$

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/b1ddbda6-4372-491f-9d99-69b8bd7aa884" style="width:66%;">
</div>


```rust
pragma circom 2.1.4;

template IsZero () {
    signal input in;
    signal output out;

    signal inv <-- in == 0? 0: 1/in;

    // 1st hint: 注意要先给out赋值，才能用第二个约束
    out <== -in * inv + 1;
    // 2st hint
    0 === in * out;
}

template Main () {
    signal input in;
    signal output o;

    component iz = IsZero();
    iz.in <== in;

    o <== iz.out;
}

component main = Main();

/* INPUT = {
    "in": "0"
} */
```

> 附：利用判0方法实现判断两数是否相等

```rust
pragma circom 2.1.4;

template IsEqual () {
    signal input in[2];
    signal output out;

    component iz = IsZero();

    iz.in <== in[0] - in[1];

    out <== iz.out;
}

template Main () {
    signal input in[2];
    signal output o;

    component ie = IsEqual();
    (ie.in[0], ie.in[1])  <== (in[0], in[1]);

    o <== ie.out;
}

component main = Main();

/* INPUT = {
    "in": ["1", "2"]
} */
```

### ②选择

`let y = s ? (a + b) : (a * b);`

由于算术电路的丰富性，需对 𝑠 进行约束检查，然后利用一个二进制位 𝑠 作为计算有效性的选择开关。

1. $𝑠 × (1 − 𝑠) = 0$
2. $𝑦 = 𝑠 × 𝑎 + 𝑏 + 1 − 𝑠 × (𝑎 ⋅ 𝑏)$

> 附：输出out应该等于in[index], 如果 index 越界（不在 [0, nChoices) 中），out 应该是 0

```rust
pragma circom 2.1.4;

// 求和
template Sum(n) {
    signal input in[n];
    signal output out;

    signal sums[n];

    sums[0] <== in[0];

    for (var i = 1; i < n; i++) {
        sums[i] <== sums[i-1] + in[i];
    }

    out <== sums[n-1];
}


template Select (nChoices) {
    signal input in[nChoices];
    signal input index;
    signal output out;

    // index 越界（不在 [0, nChoices) 中），out 应该是 0
    // 254 < 256 = 2^8
    component lt = LessThan(8);
    lt.in[0] <== index;
    lt.in[1] <== nChoices;
    lt.out === 1;

    component sm = Sum(nChoices);
    component ie[nChoices];

    for (var i = 0; i < nChoices; i++) {
        ie[i] = IsEqual();
        ie[i].in[0] <== i;
        ie[i].in[1] <== index;

        // 约束：对应index
        sm.in[i] <== ie[i].out * in[i];
    }

    // 输出out应该等于in[index]
    // 期望返回：0 + 0 ... + item
    out <== sm.out;
}

template Main () {
    signal input in[2];
    signal input index;
    signal output o;

    component st = Select(2);
    (st.in[0], st.in[1], st.index) <== (in[0], in[1], index);

    o <== st.out;
}

component main = Main();

/* INPUT = {
    "in": ["6", "7"],
    "index" : 0
} */
```

### ③二进制化

由于算术电路的丰富性，输入均为有限域 $F$ 上的数字，将其转换为二进制表示，在很多方面（比如比较大小）都有很重要的作用。与传统思路不同地方在于，将数字转化为二进制的过程，实际上是利用 $hint$ 对已经转化好的数字做约束验证的过程。

`5 -> 101`

1. $𝑜𝑢𝑡_1 × (𝑜𝑢𝑡_1 − 1) = 0$
2. $𝑜𝑢𝑡_2 × (𝑜𝑢𝑡_2 − 1) = 0$
3. $𝑜𝑢𝑡_3 × (𝑜𝑢𝑡_3 − 1) = 0$
4. $𝑜𝑢𝑡_1 × 2^0 + 𝑜𝑢𝑡_2 × 2^1 + 𝑜𝑢𝑡_3 × 2^2 = 𝑖𝑛𝑝𝑢𝑡$

```rust
pragma circom 2.1.4;

template Num2Bit (nBits) {
    signal input in;
    signal output b[nBits];

    var acc;
    for (var i = 0; i < nBits; i++) {
        b[i] <-- in \ (2 ** i) % 2;
        // 1st hint
        0 === b[i] * (1 - b[i]);
        // binary -> decimal
        acc += b[i] * (2 ** i);
    }
    
    // 2rd hint
    in === acc;
}

template Main () {
    signal input in;
    signal output o;

    component n2b = Num2Bit(4);
    n2b.in <== in;

    o <== n2b.b[0];
}

component main = Main();

/* INPUT = {
    "in": "11"
} */
```

### ④比较

`let y = s1 > s2 ? 1 : 0;`

朴素想法是将两个数字相减，将结果二进制化后，根据符号位进行判断。但由于**数字均为二进制[群](https://zhuanlan.zhihu.com/p/524518825)元素，没有负数**（我们定义算术电路上的计算都是在素数有限域上的，计算结果如果为负数需要取模，进而变为正数，**有限域上无符号位的定义**），因此我们需要将数字加上最大值，然后二进制并取最高位，通过最高位来验证。

1. $𝑦 = 𝑠_1 + 2^𝑛 − 𝑠_2$ 
2. 𝑦 二进制化取**最高位**

> 其中，n需满足 $2^n$ 大于两个参数任一个，注意n是有最大值限制(素数域)

例如：输入分别为3和4， $𝑛 = 3，𝑦 = 3 + 2^3 − 4 = 7$ ，转为二进制： $7 = (0111)_2$ ，最高位为0

```rust
pragma circom 2.1.4;

template LessThan (n) {
    signal input in[2];
    signal output out;

    // 输入信号最多为252, 2^k -1
    assert(n < 252);

    // 取n
    component n2b = Num2Bit(n+1);
    n2b.in <== in[0] + (1<<n) - in[1];

    // 二进制取最高位, 1-的目的为 < 比较，本身为 > 比较
    out <== 1 - n2b.b[n];
}

template Main () {
    signal input in[2];
    signal output o;

    component lt = LessThan(2);
    (lt.in[0], lt.in[1]) <== (in[0], in[1]);

    o <== lt.out;
}

component main = Main();

/* INPUT = {
    "in": ["0", "1"]
} */
```

### ⑤循环

由于算术电路的固定性，电路只能设计为支持最大输入数量，根据实际输入数量的不同，利用选择器技术将部分计算功能关闭，以达到不同数量的循环功能。

```rust
for (let i = 0; i < N; i++) { 
	y += 1;
}
```

利用比较方法，为临时变量 $𝑠$ 赋值，后利用选择方法，分别启用循环中的计算，或恒等原值（即未启用）。

1. $𝑠 = 1, 𝑖 < 𝑛$，否则 $s = 0$
2. $𝑦 = 𝑠 × (𝑦 + 1) + (1 − 𝑠) × (𝑦)$

### ⑥交换

通过一个交换标识 $𝑠$ 来标记是否要交换两个输入

```rust
if (s) { 
	output1 = input2; 
	output2 = input1; 
} else { 
	output1 = input1; 
	output2 = input2; 
}
```

1. $𝑜𝑢𝑡𝑝𝑢𝑡_1 = (𝑖𝑛𝑝𝑢𝑡_2 − 𝑖𝑛𝑝𝑢𝑡_1) × 𝑠 + 𝑖𝑛𝑝𝑢𝑡_1$ 
2. $𝑜𝑢𝑡𝑝𝑢𝑡_2 = (𝑖𝑛𝑝𝑢𝑡_1 − 𝑖𝑛𝑝𝑢𝑡_2) × 𝑠 + 𝑖𝑛𝑝𝑢𝑡_2$

### ⑦逻辑

```rust
let y = a & b;
let y = !a;
let y = a | b;
let y = a ^ b;
```

逻辑运算可以通过简单的数学运算获得。 另外还需要使用类似于 $a × (1 − 𝑎) = 0$ 的方式检查二进制约束。

1. $y = 𝑎 × 𝑏$ 
2. $y = 1 − 𝑎$ 
3. $y = 1 − (1 − 𝑎) × (1 − 𝑏)$ 
4. $y = (𝑎 + 𝑏) − 2 ⋅ 𝑎 × 𝑏$

### ⑧排序（冒泡排序）

```rust
for (let i = 0; i <= array.length - 1; i++) {
  for (let j = 0; j < (array.length - i - 1); j++) {
    if (array[j] > array[j + 1]) {
    	swap(array[j], array[j + 1])
    }
  }
}
```

在算术电路上做排序，可以借用**排序网络**的概念，利用多个比较器形成排序网络进行排序。

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/44ffeacd-e697-47ba-abbc-6b4c941c67fd" style="width:60%;">
</div>


## b.复杂实用电路

### ①简单签名验证

1. **KeyGen → (sk, pk)** ： 选择一个随机密钥 sk 和对应的公钥 pk
2. **Sign(m, sk) → s** ： 给定消息 m 和密钥 sk，输出签名 s
3. **Verify(m, s, pk) → 1/0**  ：给定消息 m、签名 s 和公钥 pk，验证签名是否有效

```rust
pragma circom 2.1.4;

include "circomlib/poseidon.circom";

template Sk2Pk () {
    signal input sk;
    signal output pk;

    component p;
    // 用到Poseidon组建，包含哈希过程
    p = Poseidon(1);
    p.inputs[0] <== sk;
    pk <== p.out;
}

template Sign() {
    signal input sk;
    signal input pk;

    // 生成证据时，作为public参数输入一定要有对应message（即这里产生了关联，保证了有效）
    // 但实际没有约束message(电路中即多项式中是存在的，也作为了zkp的一部分)，实际prover是对pk <> sk的证明
    // verifier 只在乎计算过程, 不关心输入吧
    signal input m;

    component s2p;
    s2p = Sk2Pk();
    s2p.sk <== sk;
    s2p.pk === pk;
}

component main = Sign();

// 其中pk通过 `Sk2Pk()` 生成
/* INPUT = {
    "sk": "6",
    "pk": "4204312525841135841975512941763794313765175850880841168060295322266705003157",
    "m": "777"
} */
```

### ②简单群签名验证

1. **KeyGen → (ski, pki)** ： 为组中的每个成员选择一组随机的秘密密钥 sk 和相应的公钥 pk 
2. **GroupSign(m, ski, G) → s** ：给定消息 m 和密钥，输出组签名 s 
3. **GroupVerify(m, s, G) → 1/0** ： 给定消息 m、组签名 s 和组 G，验证签名是否来自组

```rust
template GroupSign(n) {
    signal input sk;
    signal input pk[n];
    // 注意实际没用到的signal，建议建立自约束，防止被优化掉
    signal input m;

    component s2p;
    s2p = Sk2Pk();
    s2p.sk <== sk;
    
    signal zoreChecker[n + 1];
    // 多个乘法, 检验是否等于组内某一个
    zoreChecker[0] <== 1;
    for (var i = 0; i < n; i++) {
        zoreChecker[i + 1] <== zoreChecker[i] * (s2p.pk - pk[i]);
    }

    zoreChecker[n] === 0;
}

component main = GroupSign(2);

// 其中pk通过 `Sk2Pk()` 生成
/* INPUT = {
    "sk": "6",
    "pk": ["4204312525841135841975512941763794313765175850880841168060295322266705003157", "7"],
    "m": "777"
} */
```

### ③merkle树验证

```rust
template Reorder () {
    signal input in[2];
    // switch
    signal input s;
    signal output out[2];

    // 约束为比特位
    s * (s - 1) === 0;

    // s == 0: out[0] == in[0]
    // s == 1: out[0] == in[1]
    // out[1]同理
    out[0] <== in[0] + s * (-in[0] + in[1]);
    out[1] <== in[1] + s * (-in[1] + in[0]);
}

template MerkleVerify (n) {
    signal input root;
    signal input leaf;
    signal input siblings[n];
    signal input pathIndex[n];

    signal hashes[n + 1];

    component reorder[n];
    component poseidon[n];
    hashes[0] <== leaf;

    for (var i = 0; i < n; i++) {
        // 1. 排序
        reorder[i] = Reorder();
        reorder[i].in[0] <== hashes[i];
        reorder[i].in[1] <== siblings[i];
        reorder[i].s <== pathIndex[i];

        // 2.哈希
        poseidon[i] = Poseidon(2);
        poseidon[i].inputs[0] <== reorder[i].out[0];
        poseidon[i].inputs[1] <== reorder[i].out[1];
        hashes[i + 1] <== poseidon[i].out;
    }

    // 约束验证
    hashes[n] === root;
}

component main { public [root, leaf] } = MerkleVerify(15);

/* INPUT = {
    "root": "12890874683796057475982638126021753466203617277177808903147539631297044918772",
    "leaf": "1355224352695827483975080807178260403365748530407",
    "siblings": [
        "1",
        "217234377348884654691879377518794323857294947151490278790710809376325639809",
        "18624361856574916496058203820366795950790078780687078257641649903530959943449",
        "19831903348221211061287449275113949495274937755341117892716020320428427983768",
        "5101361658164783800162950277964947086522384365207151283079909745362546177817",
        "11552819453851113656956689238827707323483753486799384854128595967739676085386",
        "10483540708739576660440356112223782712680507694971046950485797346645134034053",
        "7389929564247907165221817742923803467566552273918071630442219344496852141897",
        "6373467404037422198696850591961270197948259393735756505350173302460761391561",
        "14340012938942512497418634250250812329499499250184704496617019030530171289909",
        "10566235887680695760439252521824446945750533956882759130656396012316506290852",
        "14058207238811178801861080665931986752520779251556785412233046706263822020051",
        "1841804857146338876502603211473795482567574429038948082406470282797710112230",
        "6068974671277751946941356330314625335924522973707504316217201913831393258319",
        "10344803844228993379415834281058662700959138333457605334309913075063427817480"
    ],
    "pathIndex": [
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1",
        "1"
    ]
} */
```

## c.Circom

### ①理论与语法

[Circom](https://docs.circom.io/)是一个底层用rust实现的开源编译器，它可以编译用circom语言实现的电路circuit。它将circuit编译的结果以contraints的形式输出，这些constraints能被用于计算相应生成逻辑的proof。

> 在线编译器：https://zkrepl.dev/

**circom在整个过程中的逻辑关系：**

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/61c77b4b-b310-4d62-b90b-33506c907672" style="width:70%;">
</div>

**各文件的功用：**

- Circuit.circom：程序员编写的电路代码
- Input.json：输入，如public input
- PoT.ptau：proof tau，根据约束的**随机数**文件，约束越多，需要匹配的随机数的消耗越大
- Circuit.wasm：webassembly证明器
- Proving Key (.zkey)
- Verification Key (.vkey)
- Verifier.sol：solidity验证器（也可以是node.js服务端）



1. **circom** 编写电路代码

   > 从编程角度，仅仅这一步是需要开发者做的，后面两步均为编译器生成。其中写circom代码就像是在证明器里做验证器工作。
   >
   > 对比JS等语言的独特设计：1.中间过程值暴露 2.暴露出约束

2. 生成 **wtns+r1cs** 约束文件

3. 套用框架(groth16/plonk)生成**ZKP证据**与**solidity验证合约**

   > 实际ZK耗时点在这一步。PLONK模型比groth16的数学难题更简单，但是难题的难度不减，实际上在数学中这被认为是更好的选择。



### ②代码实践

**一定注意**，circom中有分**Konwn和[Unkonwns域](https://docs.circom.io/circom-language/circom-insight/unknowns/#control-flow)**，代码里面的 for/if 是生成电路用的，而非电路里的具体约束逻辑，电路逻辑中是没有for/if等一下的逻辑的，需要自己设计（对应电路基础中的案例）。

> 零知识证明电路二进制化设计，circom实现：

```c++
pragma circom 2.1.4;

include "circomlib/poseidon.circom";
// include "https://github.com/0xPARC/circom-secp256k1/blob/master/circuits/bigint.circom";

template Num2Bits (nBits) {
    signal input in;
    signal output b[nBits];

    var acc;
    for (var i = 0; i < nBits; i++) {
        // 注意“\”整除
        b[i] <-- (in \ 2 ** i) % 2;
        // 约束为0或1
        0 === b[i] * (1 - b[i]);
        // 累加器
        acc += b[i] * (2 ** i);
    }

    // 真正的电路就是“===”，作约束
    in === acc;
}

template Main () {
    signal input in;
    signal output out;

    // example
    component n2b = Num2Bits(4);
    n2b.in <== in;

    // 导出第0位
    // <== 等于 <-- + ===
    out <== n2b.b[0];
}

component main = Main();

/* INPUT = {
    "in": "11"
} */
```

### ③哈希算法比较

主流哈希算法的效率比较：

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/76a7d537-d73d-4337-b82c-c463965273c9" style="width:50%;">
</div>

其次，基于区块链的ZK实现，应选择**链上计算**和**电路计算**都高效的（对应到EVM上就是同时要考虑Gas Cost），更多请查看[这篇](https://ethresear.ch/t/gas-and-circuit-constraint-benchmarks-of-binary-and-quinary-incremental-merkle-trees-using-the-poseidon-hash-function/7446)。

| [Rank(Best to Worst)](https://twitter.com/RAILGUN_Project/status/1363686166734675968?s=20) |      ETH Gas Cost       |  ZK Circuit Constraint  |
| :----------------------------------------------------------: | :---------------------: | :---------------------: |
|                              1                               |        Keccak256        |  Poseidon T6(Quinary)   |
|                            **2**                             |         SHA256          | **Poseidon T3(Binary)** |
|                            **3**                             | **Poseidon T3(Binary)** |       MiMC Sponge       |
|                              4                               |  Poseidon T6(Quinary)   |        Keccak256        |
|                              5                               |       MiMC Sponge       |         SHA256          |


<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/d80f2ce3-abb9-42bf-a6bb-fe1adada283b" style="width:30%;">
</div>


## d.应用ZK架构

### ①递归和组合

**递归**：递归使得复杂计算的证明可以**并行化**

**证明组合**：将来自**不同证明系统的子协议**嵌在一起。

1）递归技术：

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/2272aa4c-44c1-419e-8f36-932b017bbbc5" style="width:70%;">
</div>

**例：[IVC(Incrementally Verifiable  Computation)完全递归](https://www.cs.purdue.edu/homes/pvaliant/uniqueCS.pdf)**

> 如Plonky2、Nova算法，利用递归组合技术

策略：将 $𝑧_𝑛 = 𝐹^{(𝑛)} (𝑧_0; 𝑤_0, … , 𝑤_{𝑛−1})$ 分解为 $𝐹$ 的递归应用

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/7f535c8e-fb68-465f-85ac-782dc76e82bc" style="width:70%;">
</div>

**应用**：

- **每个区块可以常数时间内验证的区块链**：证据 $π$ 的递归引用前一次的（除了genesis），每次只需要验证最后一个证据可证明所有历史

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/bf6607f6-2442-4d9b-adbd-f82eff361e56" style="width:70%;">
</div>

- **可验证延迟函数[BBBF18]**：VDF做不到并行化，要一层一层算过去，O（n）。利用ZK可以做到数据计算+前面的证据是正确的，通过验证并行化，递归所有证据（这些证据不需要相同算法/框架算出，这里也是ZKML的一个应用点）后得到一个证据，只验证一个证据O(1)就可以验证全部。**目的是降低验证速度**（但计算复杂性可能会提升）

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/3870b408-634c-4284-b1e9-44ccc41ace9f" style="width:70%;">
</div>

2）组合技术：递归验证

|             | 快速证明器（证明的时间复杂度） | 小证据/快速验证器（证据空间复杂度） |
| :---------: | :----------------------------: | :---------------------------------: |
|  **STARK**  |            ✅ $O(n)$            |             ❌ $O(logn)$             |
| **Groth16** |            ❌ $O(n)$            |              ✅ $O(1)$               |

1. 利用STARK的快速证明特性，设计STARK电路

2. 生成中间大证据 $π_{STARK}$

3. 在SNARK (如上述 Groth16) 电路中，实现STARK验证器

   > 这里SNARK虽然是慢速证明器，但是要证明的内容是 $O(logn)$ 而非 $O(n)$ ，复杂度是降下来的

4. 生成最终小证据 $π_{SNARK}$

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/0710f799-89d3-44c8-a9a9-5b63e732afec" style="width:70%;">
</div>

**应用**：

- **可链接的提交与证明**：zkSNARKs Portfolio(proof gadgets)

  > legoSNARKs[CFQ19]

  类似Pedersen承诺(2.2.a章节)的 ”同态性“，期望通过数学构造的方式，根据不同逻辑复杂性，选择不同的证明器，然后高效的把所有证据合在一起达成一个证据。

  <div align=center>
  <img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/8cb7dd07-b29a-444a-88f8-791b82984f19" style="width:45%;">
  </div>


### ②[零知识数据市场](https://github.com/nulven/EthDataMarketplace)

> 应用实例：Dark Forest 创建的 [NightMarket](https://nightmart.xyz/) 

传统数据市场模式：Escrow第三方托管，而在依赖区块链的透明网络的市场交易过程中，如何做到不对其他人透露信息的情况下完成交易？解决方案（简单实现）：

> 现在有一个场景，拥有资金的Bob，想从拥有信息的Alice手中购买信息

1. $Alice$ 使用 $Bob$ 的公钥加密数据并加以发布。
2. 同时，$Alice$ 还需要发布 $zkSNARK$ 证据，证明该密文是使用 $Bob$ 的公钥正确加密的数据。 

只有当 $zkSNARK$ 证据被验证后，智能合约才会向 $Alice$ 释放资金。过程中涉及到的参数如下所示

公共输入：

- 买方(Bob)公钥 $pk$ 
- 密文 $c$ 
- 承诺 $h$ 

私密输入： 

- 隐私数据 $s$ 

证明：

- $Hash(s) = h$ 
- $Enc(s, pk) = c$

### ③[ZKML](https://github.com/zkonduit/ezkl)

神经网络是一个函数，ZKSNARKs本身也是一个函数，可以将其放到ZKSNARKs中。

**应用实例一**：可验证计算

> 场景：希望通过 LLM(大语言模型) 来进行案件审批/专家建议与推断，而模型本身并非是公开的 (e.g.CloseAI）
>
> 由谁来运行模型？如何验证结果的正确性？

在任何案例开始之前：LLM 承诺模型 $（Model_Commit) = C$

公共输入： 

- 输入 $x$ 
- 模型声明的输出 $y$ 
- 模型承诺 $C$ 

隐私输入： 

- 模型 $M$ 

证明： 

- $M(x) = y$ 
- $commit(M) = C$

**应用实例二**：零知识生物识别 ZK Biometrics （e.g. [Worldcoin](https://worldcoin.org/)）

> 场景：生物识别认证目前只有在大型机构存储，如何在不泄露个人隐私的情况下，实现公共生物识别数据认证？

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/90cb5ea5-c0f9-4425-a66e-5c0e1b9b3cae" style="width:60%;">
</div>

### ④未来方向

- **基准测试**
  - 多项式承诺基准测试：不同框架下对方法实现不完善，比如在hash函数方面
  - ZKP编译器优化：如circom检查比较少
- **库、标准、开发工具**
  - 如legoSNARKs在做的
  - 能否为递归/组合策略定义高级API
- **安全分析**
  - 理论安全，公式/代码角度的安全需要保证
  - 递归组合的安全性：需要更强的知识提取器（必须在每个递归步骤中都成功）
  - 利用统一框架简化分析

## e.应用ZK实例

### ①Group Membership

> Nullifiers：无效化参数，使用户无法匿名地执行两次相同操作（注意并不和投票人关联，即只知道投过票不知道谁投的）

<div align=center>
<img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/c6934baa-cb53-496c-adb9-427aba90e27a" style="width:30%;">
</div>

- **无Nullifiers**

  > Spec1：[zkMessage](https://zkmessage.xyz/)：通过每条消息附带的 zk 成员证明，代表一个群组发布消息。

  - step1：注册阶段

    User 生成 hash(9j...46)=75...gk 其中的**秘密字符串**：9j...46，公开的是哈希，存入Database

  - step2：发消息阶段，message="hi"

    > 具体circom实现见算术电路应用的群签名部分

    **prove**(9j...46, [75...gk, ...(database所有成员的公开变量hash)], “hi”)

  *弊端：需要将secret字段传入，即使是隐私变量*

  > Spec2 ：[heyanon](https://twitter.com/heyanonxyz)：公私钥对的签名代替存入secret

- **有Nullifiers**

  > Spec3 ：[Semaphore](https://semaphore.appliedzkp.org/)：证明哈希列表中的成员身份，包含id commitment, nullifier, [trapdoor](陷门，是私有的。陷门函数：正向计算是很容易的，但若要有效的执行反向计算则必须要知道一些secret/key/knowledge/trapdoor)

  - step1：注册阶段

    user 以哈希身份加入： hash(nullifier, trapdoor)，nullifier, trapdoor 是私有的，存入数据库的只有hash本身，数据库内可用merkle tree 构建

  - step2：发消息阶段(对问题 “Does pineapple belong on pizza?” 投票)

    **prove**(merkle_root=0x59..., “Yes” , “Does...”, nullifier1, trapdoor1, merkle_path)

  > Spec4 ：[tornado cash](https://ipfs.io/ipns/tornadocash.eth/)：通过每次提款附加一个zk成员身份证明和一个nullifier，实现向匿名账户发送资金
  >
  > 混币器Coin Mixer：同Mixer币需要等额，无区分度；不同额度不同Mixer（有的交易所会全链路检查，不接受tornado cash出的钱，不过现在tornado cash提供一个证据，帮助用户出据全链路证明合法）

  - 存款阶段

    User 发送 1 eth 并 hash(**nullifier** | **secret**，入merkle tree

    > secret：辅助提款信息，存款人私下生成secret时已经定住了，但不一定与user address相同

  - 提款阶段

    <div align=center>
    <img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/3bcd4896-c777-489b-9783-2efba3b93011" style="width:30%;">
    </div>

    user链下操作：**prove**(merkle_root=0x59..., **recipient_pk**=0x7ab89.., **nullifier**, **secret**, merkle_path1)

    *注意：要保证匿名传入与接受不对称，需要传入接受者地址，该地址是融入zk电路中的，frontrun没有用*

    <div align=center>
    <img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/ae0b2e91-a420-493a-ae30-891b643f1eb1" style="width:35%;">
    </div>


  > Spec5 ：[Stealthdrop](https://stealthdrop.xyz/)：通过每次提款附加一个zk成员身份证明和一个nullifier，向无关联账户 (unlinked account)申领空投
  >
  > ***注意：存在bug**

  <div align=center>
  <img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/bb0a3257-4bba-4ee2-addc-3c22540301c7" style="width:60%;">
  </div>


  - **BUG**：如图所示，nullifier_hash = hash(signature_of_message)。而**ECDSA为非确定性签名**，会依赖一个随机数从而使**不同随机数对应不同签名**，而nullifier是签名的哈希，自然每次都是不同的nullifier，可以无限领取空投。

    > 即使是用确定性签名，也**不可行**。
    >
    > 确定性签名：**通过一个 “数K“ 代替了随机数**，且通过数学难题使得无法通过签名反推出私钥，以保证安全性。但在此处仍不可行，因为换一个数K签名仍然有效，该数是可以更换的，故而验证方仍然无法辨别。

  - **解决方案**：传入Private Key 而非 Signature，但更好的是：基于用户私钥的确定性函数， 该函数可以仅通过用户的公钥进行验证， 并保持其匿名，例如通过secret key算出群上面的元素 $hash(message, public key)^{secret key} => DDH-VRF$ （Decisional Diffie–Hellman Verifiable Random Function ）

- **其他范例：zk-Email**

  不公开邮件内容证明我收到过某封邮件

  [DKIM](https://dmarcly.com/blog/zh-CN/dkim-faqs-frequently-asked-questions)：域名运营者 (邮件服务器) 的来操作，与用户无关。此外，这仅仅是域名层面的过滤，是不会对内容进行过滤

  <div align=center>
  <img src="https://github.com/DessertHeart/Dive-Into-Blockchain/assets/93460127/f8d06466-ac83-4e78-a15e-d2feb78f1dc1" style="width:60%;">
  </div>


  应用：

  1. **匿名KYC**：即证明收到了如Binance的有效KYC电子邮件，以证明是人类但不暴露个人隐私信息
  2. **银行余额认证**：通过中心化银行账户余额电子邮件，证明您的银行账户中有X元

  ……

### ②Dark Forest

一款MMO（大型多人在线游戏类型）游戏，也是第一个**全链游戏**，是一个以太坊智能合约，具有**无许可的互操作性**。故而有很多开发的 **“插件”** （如著名的代理hash插件 remote explorer）可以接入。

**[游戏构造](https://blog.zkga.me/df-init-circuit)：**

> 每个玩家都在一个大的二维网格上（高亮部分即有活动）

1. **玩家状态**

   - 公共状态：拥有哪些公共地址，谁拥有它们以及它们的人口数量

   - 私有状态：玩家行星的私有地址 $(x, y)$ 

   对于位置 $(x,y)$，$hash(x,y)$ 是该位置的**公共地址**，这些**坐标本身**是该位置的**私有地址**。

   当 $hash(x,y)$ < `DIFFICULTY_THRESHOLD` (该值的大小，决定了星球的大小)时，位置 $(x,y)$ 上有适合居住的星球。如果不满足，则该空间是空的，由玩家控制的单位仅存在于玩家拥有的星球上

2. **玩家行为**

   - 探索 (初始化)

     星球有两种 “资源”：人口和矿。两个参数都会缓慢增长但是有上限，拥有足够的 "矿" 可以升级星球。

     $zkProof$ : 证明我知道某坐标 $(x,y)$ ，使得

     - $hash(x, y) = planetId$ 

     - $𝑥^2 + 𝑦^2 < 𝑐𝑙𝑎𝑖𝑚𝑒𝑑𝐷𝑖𝑠𝑡^2$

       > $planetId$ ：星球坐标的 $hash$ ，即 $hash(x, y)$
       >
       > $𝑐𝑙𝑎𝑖𝑚𝑒𝑑𝐷𝑖𝑠𝑡$ ：该星球位置到坐标原点 $(0,0)$ 的距离

   - 占领 (移动)

     移动的同时可以指明携带的资源，如果携带的人员超过该星球的人口，说明可以攻占星球，但需支付一些费用，具体取决于两个星球之间的**最大距离**，星球间移动存在移动速度。

     $zkProof$ : 为了在不公开星球坐标的情况，还能证明星球的移动正确，我知道某坐标 $(x_1, y_1)$ 和 $(x_2, y_2)$ ，使得

     - $hash(x_1, y_1) = fromPlanetId$

     - $hash(x_2, y_2) = toPlanetId$

     - $x_2^2 + y_2^2 < worldRadius^2$

     - $(x_1-x_2)^2 + (y_1-y_2)^2 < distMax^2$

       > $distMax$ : 星球间最大距离
       >
       > $worldRadius$ :整个宇宙坐标的最大半径，即检查这两个星球是否 “在边界内” 
