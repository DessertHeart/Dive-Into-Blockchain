// SPDX-License-Identifier: MIT

pragma solidity 0.8.10;

// UniswapV3
import './uniswap-v3-periphery/interfaces/ISwapRouter.sol';
// UniswapV2
import './uniswap-v2-periphery-master/contracts/interfaces/IUniswapV2Router02.sol';
// AAVE
import {FlashLoanSimpleReceiverBase} from './aave-v3-core/contracts/flashloan/base/FlashLoanSimpleReceiverBase.sol';
import {IPoolAddressesProvider} from './aave-v3-core/contracts/interfaces/IPoolAddressesProvider.sol';
// Others
import {IERC20} from './aave-v3-core/contracts/dependencies/openzeppelin/contracts/IERC20.sol';
import './uniswap-v3-core/libraries/TransferHelper.sol';


contract AaveFlashswap is FlashLoanSimpleReceiverBase {

    //rinkeby address

    address private constant WETH = 0xc778417E063141139Fce010982780140Aa0cD5Ab;
    address private constant ATOKEN =0x784c47Ba17A32e9C636cf917c9034c0aD1E87d41;
    address private constant UNISWAP_V2_ROUTER =0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D;
    address private constant SWAPROUTER =0xE592427A0AEce92De3Edee1F18E0157C05861564;
    //V3池费
    uint24 private constant  poolFee = 3000;
    
    // AAVE.pool合约地址
    address private constant PERSONAL =0xAe16D9E58c63cf274E30B06Cb7c9C5367c3229C9;

    event SuccessEvent(string indexed message);
    event CatchStringError(string indexed message);
    event CatchDataError(bytes indexed data);

    // IPoolAddressesProvider?
    constructor(IPoolAddressesProvider _provider) FlashLoanSimpleReceiverBase(_provider) public {}

    // 调用swap执行闪电贷
    function flashSwapCall(address _asset, uint _amount) public {

        // POOL 来自FlashLoanReceiverBase.
        try POOL.flashLoanSimple({
            receiverAddress: address(this),
            asset: _asset,
            amount: _amount,
            params: abi.encode(uint256(_amount)),
            referralCode: 0
        }) {
            // 无利可图，回滚
            uint256 balance = IERC20(_asset).balanceOf(address(this));
            require(balance > 0, "AAVE flashloan Fail");

            emit SuccessEvent("FlashSwap Success!");
        } catch Error(string memory reason) {
            emit CatchStringError(reason);
        } catch (bytes memory data) {
            emit CatchDataError(data);
        }
        // 将套利结果从合约中提出
        TransferHelper.safeTransfer(_asset, msg.sender, balance);
    }

    // 注意，_sender为调用POOL.flashLoanSimple的地址，这里是实际上就是本合约
    function executeOperation(
        address _asset,
        uint256 _amount,
        uint256 _fees,
        address _sender,
        bytes memory _params
    ) public override returns(bool){

        // 测试用授权
        IERC20(WETH).approve(UNISWAP_V2_ROUTER, 100_000);
        IERC20(ATOKEN).approve(UNISWAP_V2_ROUTER, 100_000);
        IERC20(WETH).approve(SWAPROUTER, 100_000);
        IERC20(ATOKEN).approve(SWAPROUTER, 100_000);

        address[] memory path = new address[](2);
        path[0] = WETH;
        path[1] = ATOKEN;
        amountLoan = IERC20(WETH).balanceOf(address(this));
        //AAVE借来的WETH，V2买ATOKEN
        IUniswapV2Router02(UNISWAP_V2_ROUTER).swapExactTokensForTokens(
            amountLoan,
            0,
            path,
            address(this),
            block.timestamp+2000
        );
        // V2买到的ATOKEN
        uint256 amountToken = IERC20(ATOKEN).balanceOf(address(this));

        //V3通过ATOKEN买WETH(注意：V3购买没有deadline参数，与github上不一致)
        ISwapRouter.ExactInputSingleParams memory params = ISwapRouter
            .ExactInputSingleParams({
                tokenIn: ATOKEN,
                tokenOut: WETH,
                fee: poolFee,
                recipient: address(this),
                amountIn: amountToken,
                amountOutMinimum: 0,
                sqrtPriceLimitX96: 0
            });
        // V3获得的WETH
        uint256 amountOut = ISwapRouter(SWAPROUTER).exactInputSingle(params);

        // 还款AAVE
        uint256 amountRequired = _amount + _fees;

        // 偿还只要授权即可，AAVE会自己尝试转账
        IERC20(_asset).approve(msg.sender, amountRequired);
        // 套利
        TransferHelper.safeTransfer(_asset, _sender, amountOut - amountRequired);

        return true;
    }

}