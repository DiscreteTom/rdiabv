package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/DiscreteTom/rdiabv"
)

const timesForHTRM = 10

func runHTRM() bool {
	// open files
	dataFile, err := os.Open(dataFilename)
	if err != nil {
		panic(err)
	}
	tagFile, err := os.Open(tagFilename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(tagFile)

	// init read buffer for data file
	dataBuffer := make([]byte, chunkSize)

	// init dhdd
	fmt.Println("Initializing HTRM...")
	htrm := rdiabv.NewHTRM(timesForHTRM).
		InitBuffers(NewRawRsaBlock())

	// read & merge
	fmt.Println("Merging...")
	for i := 0; i < blockCount; i++ {
		// for data file, read a chunk
		_, err := dataFile.Read(dataBuffer)
		if err != nil {
			panic(err)
		}
		// create block
		var block = RawRsaBlock{Data: new(big.Int).SetBytes(dataBuffer)}
		// for tag file, read a line
		scanner.Scan()
		block.Tag, _ = new(big.Int).SetString(scanner.Text(), 10)
		htrm.MergeBlock(&block)
	}

	fmt.Println("Checking...")
	return htrm.CheckAllBuffers()
}
