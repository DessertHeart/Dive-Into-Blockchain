package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

/* ------------------------------CLI界面---------------------------------- */
/* -----------------------------遍历区块链---------------------------------- */


// CLI
type CLI struct {}

// CLI工具栏介绍
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address 'YOUR ADDRESS")
	fmt.Println("  createwallet")
	fmt.Println("  getbalance -address 'YOUR ADDRESS'")
	// fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain")
	fmt.Println("  send -from 'FROM' -to 'TO' -value 'AMOUNT'")
	//fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

// CLI入口方法
func (this *CLI) Run() {
	// 输入参数检查
	if len(os.Args) < 2 {
		fmt.Println("No enough parameters!")
		fmt.Println()
		this.printUsage()
		os.Exit(1)
	}

	// CreateBlockchain
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	// Send
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	fromData := sendCmd.String("from", "", "Sender address")
	toData := sendCmd.String("to", "", "Receiver address")
	valueData := sendCmd.Int("value", 0, "Amount of satoshi")

	// Print
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	// CreateWallet
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	// GetBalance
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")

	// os.Args[]存储命令行参数，[0]为程序名
	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send" :
		// 从Args[]中解析注册的flag
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {log.Panic("addBlockCmd.Parse failed")}
	case "printchain" :
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {log.Panic("printChainCmd.Parse failed")}
	case "getbalance" :
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		os.Exit(1)
	}
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		this.createBlockchain(*createBlockchainAddress)
	}

	// flag.Parse()是否已被调用
	if sendCmd.Parsed() {
		if *fromData == "" || *toData == "" || *valueData <= 0 {
			// 解析失败
			sendCmd.Usage()
			os.Exit(1)
		}
		// 执行添加区块命令
		this.Send(*fromData, *toData, *valueData)
	}

	if printChainCmd.Parsed() {
		// 执行打印命令
		this.PrintChain()
	}

	if createWalletCmd.Parsed() {
		// 执行创建新钱包命令
		this.createWallet()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		//cli.addBlock(*addBlockData)
		this.getBalance(*getBalanceAddress)
	}
}