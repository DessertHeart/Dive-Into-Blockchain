# :zap:Flashloan

闪电贷，顾名思义，就是利用交易的原子性，在一个区块时间或者一笔交易内(快如闪电的耗时:smirk:)完成贷款与还款的操作。<br>
闪电贷在2020的DeFi Summer扮演了十分重要的角色，为区块链的金融业务提供了免抵押借款服务，是DeFi世界的一款利器。<br>
闪电贷的概念最早是由Marble协议提出来的，第一笔闪电贷操作来自于Aave协议。

## Resolution

实现功能：
- Step1：从AAVE借贷来 WETH，Uniswap-V2 swap WETH => AT Token
- Step2：通过 Uniswap-V3 swap AT Token => WETH
- Step3：V3 拿到的WETH还款给 AAVE，如果有利润可以实现差价套利

![image](https://user-images.githubusercontent.com/93460127/210203066-ed8768b8-0926-4894-8ff7-a1664d3e34b2.png)


