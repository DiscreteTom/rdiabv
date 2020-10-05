package main

import (
	"bufio"
	"math/big"
	"os"
)

// check every block
func runOneByOne() bool {
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

	// check each block
	for i := 0; i < blockCount; i++ {
		// for data file, read a chunk
		_, err := dataFile.Read(dataBuffer)
		if err != nil {
			panic(err)
		}
		var data = new(big.Int).SetBytes(dataBuffer)
		scanner.Scan()
		tagFromFile, _ := new(big.Int).SetString(scanner.Text(), 10)
		if data.Cmp(rawRsa.RawDecrypt(tagFromFile)) != 0 {
			return false
		}
	}
	return true
}
