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
}

func (e *Executor) PrintChain() {
	reader := e.bc.NewReader()

	for len(reader.current) != 0 {
		b, err := reader.Next()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Data - %s, Hash %x, Perev %x \n", b.Data, b.Hash, b.PrevBlockHash)

		if len(b.PrevBlockHash) == 0 {
			fmt.Println("-- EOC")
		}
	}
}

func (e Executor) run() {
	e.validateArgs()
	addBlock := flag.NewFlagSet("addblock", flag.ExitOnError)
	printchain := flag.NewFlagSet("printchain", flag.ExitOnError)

	addData := addBlock.String("data", "", "Data to be added (String)")

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
	default:
		e.printUsage()
		os.Exit(1)
	}

	if addBlock.Parsed() {
		if *addData == "" {
			e.printUsage()
			os.Exit(1)
		}
		e.bc.AddBlock(*addData)
	}

	if printchain.Parsed() {
		e.PrintChain()
	}

}

func NewCommandExecutor(bc *Blockchain) *Executor {
	return &Executor{bc}
}
