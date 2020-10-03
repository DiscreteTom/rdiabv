package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/DiscreteTom/rdiabv"
)

const rsaKeyBits = 4096
const chunkSize = 256 // Byte count of a block's data field.
const blockCount = 1024
const timesForHTRM = 10
const dataFilename = "data.bin"
const tagFilename = "tag.txt"

var rawRsa = NewRawRsa(rand.Reader, rsaKeyBits)

func main() {
	// generate test data
	fmt.Println("Generating files...")
	fileGen(dataFilename, tagFilename, chunkSize, blockCount)

	trackDuration("DHDD", runDHDD)
	trackDuration("HTRM", runHTRM)
	trackDuration("One by one", runOneByOne)
}

func trackDuration(name string, f func() bool) {
	fmt.Println("Running " + name + ":")
	start := time.Now()
	fmt.Printf("Result: %v\n", f())
	fmt.Printf("%v: %v\n", name, time.Since(start))
}

func runDHDD() bool {
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
	fmt.Println("Initializing DHDD...")
	dhdd := rdiabv.NewDHDD(blockCount, time.Now().UnixNano()).
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
		dhdd.MergeBlock(i, &block)
	}

	fmt.Println("Checking...")
	return dhdd.CheckAllBuffers()
}

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
		if data.Cmp(rawRsa.Decrypt(tagFromFile)) != 0 {
			return false
		}
	}
	return true
}
