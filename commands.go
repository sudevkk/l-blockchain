package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Executor struct {
	bc *Blockchain
}

func (cli *Executor) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *Executor) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
	fmt.Println("  getbalance -address (Get Balance Amount of Address)")
}

func (e *Executor) PrintChain() {
	reader := e.bc.NewReader()

	for len(reader.current) != 0 {
		b, err := reader.Next()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Data - %+v, Hash %x, Perev %x \n", b.Transactions, b.Hash, b.PrevBlockHash)

		if len(b.PrevBlockHash) == 0 {
			fmt.Println("-- EOC")
		}
	}
}

func (cli Executor) getBalanceOf(address string) {
	balance := cli.bc.getBalance(address)
	fmt.Printf("Balance of %s - %d", address, balance)
}

func (cli Executor) transferAmount(from string, to string, amount int) {
	cli.bc.transfer(from, to, amount)
}

func (e Executor) run() {
	e.validateArgs()
	addBlock := flag.NewFlagSet("addblock", flag.ExitOnError)
	printchain := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalance := flag.NewFlagSet("getbalance", flag.ExitOnError)
	send := flag.NewFlagSet("send", flag.ExitOnError)

	addData := addBlock.String("data", "", "Data to be added (String)")
	address := getBalance.String("address", "", "Address (String)")
	from := send.String("from", "", "From Address (String)")
	to := send.String("to", "", "To Address (String)")
	amount := send.Int("Amount", 0, "Amount to be sent (String)")

	switch os.Args[1] {
	case "addblock":
		fmt.Println("Add Block Command")
		err := addBlock.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		fmt.Println("Printchain Command")
		err := printchain.Parse(os.Args[1:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		fmt.Println("Getbalance Command")
		err := getBalance.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		fmt.Println("Send Command")
		err := send.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		e.printUsage()
		os.Exit(1)
	}

	if addBlock.Parsed() {
		if *addData == "" {
			e.printUsage()
			os.Exit(1)
		}
		e.bc.AddBlock([]*Transaction{})
	}

	if printchain.Parsed() {
		e.PrintChain()
	}

	if getBalance.Parsed() {
		if *address == "" {
			e.printUsage()
			os.Exit(1)
		}
		e.getBalanceOf(*address)
	}
	if send.Parsed() {
		if *from == "" || *to == "" || *amount == 0 {
			e.printUsage()
			os.Exit(1)
		}
		e.transferAmount(*from, *to, *amount)
	}

}

func NewCommandExecutor(bc *Blockchain) *Executor {
	return &Executor{bc}
}
