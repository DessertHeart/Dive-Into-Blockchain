let { ethers } = require("hardhat");

async function main() {
    let [owner] = await ethers.getSigners();

    let provider = "0xBA6378f1c1D046e9EB0F538560BA7558546edF3C";
    let weth = "0xc778417E063141139Fce010982780140Aa0cD5Ab";
    
    let FlashLoanAave = await ethers.getContractFactory("FlashLoanAave");
    let flashLoanAave = await FlashLoanAave.deploy(provider,{ gasLimit: 8000000 });
    await flashLoanAave.deployed();
    console.log("flashLoanAave:" + flashLoanAave.address);

    let loanAmount = ethers.utils.parseUnits("1", 17); //借0.1个以太
    await flashLoanAave.flashSwap(weth, loanAmount, { gasLimit: 8000000 });
    console.log("发起AAVE闪电贷");
}   

 
main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });
 