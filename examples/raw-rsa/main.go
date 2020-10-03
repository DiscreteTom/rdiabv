package main

import (
	"crypto/rand"
	"fmt"
	"time"
)

const rsaKeyBits = 4096
const chunkSize = 256 // Byte count of a block's data field.
const blockCount = 1024
const dataFilename = "data.bin"
const tagFilename = "tag.txt"

var rawRsa = NewRawRsa(rand.Reader, rsaKeyBits)

func main() {
	// generate test data
	fmt.Println("Generating files...")
	fileGen(dataFilename, tagFilename, chunkSize, blockCount)

	// performance test
	trackDuration("DHDD", runDHDD)
	trackDuration("HTRM", runHTRM)
	trackDuration("One by one", runOneByOne)
}

// call f and track duration
func trackDuration(name string, f func() bool) {
	fmt.Println("Running " + name + ":")
	start := time.Now()
	fmt.Printf("Result: %v\n", f())
	fmt.Printf("%v: %v\n", name, time.Since(start))
}
