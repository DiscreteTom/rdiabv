package main

import (
	"math/big"
	"math/rand"
	"os"
)

// fileGen will generate random data file(binary) and it's tag file(text).
func fileGen(dataFilename, tagFilename string, chunkSize, blockCount int) {
	// create data file
	dataFile, err := os.OpenFile(dataFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()
	// create tag file
	tagFile, err := os.OpenFile(tagFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer tagFile.Close()

	// write files
	data := make([]byte, chunkSize)
	for i := 0; i < blockCount; i++ {
		// data file consists of random bytes
		rand.Read(data)
		if _, err = dataFile.Write(data); err != nil {
			panic(err)
		}
		// each line in tag file is a tag
		if _, err = tagFile.WriteString(rawRsa.Encrypt(new(big.Int).SetBytes(data)).String() + "\n"); err != nil {
			panic(err)
		}
	}
}
