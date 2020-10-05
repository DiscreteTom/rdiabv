package main

import (
	"fmt"
	"time"
)

func main() {
	const rsaKeyBits = 4096
	const chunkSize = 256 // Byte count of a block's data field.
	const blockCount = 1024
	const dataFilename = "data.bin"
	const tagFilename = "tag.txt"
	const keyFilename = "key.pem"

	// generate key files and test data
	fmt.Println("Generating files...")
	fm := NewFiles(keyFilename, rsaKeyBits, dataFilename, tagFilename, chunkSize, blockCount)

	// or you can load existing files
	// fmt.Println("Using existing files.")
	// fm := LoadFiles(keyFilename, dataFilename, tagFilename, chunkSize, blockCount)

	// performance test
	trackDuration("DHDD", runDHDD, fm)
	trackDuration("HTRM", runHTRM, fm)
	trackDuration("One by one", runOneByOne, fm)
}

// call f and track duration
func trackDuration(name string, f func(*FileManager) bool, fm *FileManager) {
	fmt.Println("Running " + name + ":")
	start := time.Now()
	fmt.Printf("Result: %v\n", f(fm))
	fmt.Printf("%v: %v\n", name, time.Since(start))
}
