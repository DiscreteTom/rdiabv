package main

import (
	"fmt"
	"time"

	"github.com/DiscreteTom/rdiabv"
)

func runDHDD(fm *FileManager) bool {
	fm.StartSession()
	defer fm.EndSession()

	// init dhdd
	fmt.Println("Initializing DHDD...")
	dhdd := rdiabv.NewDHDD(fm.blockCount, time.Now().UnixNano()).
		InitBuffers(NewRawRsaBlock(fm.rr))

	// read & merge
	fmt.Println("Merging...")
	for i := 0; i < fm.blockCount; i++ {
		// create block
		var block = RawRsaBlock{Data: fm.NextBlockData(), Tag: fm.NextBlockTag(), Key: fm.rr}
		dhdd.MergeBlock(i, &block)
	}

	fmt.Println("Checking...")
	return dhdd.CheckAllBuffers()
}

func parallelRunDHDD(fm *FileManager) bool {
	fm.StartSession()
	defer fm.EndSession()

	// init dhdd
	fmt.Println("Parallel initializing DHDD...")
	dhdd := rdiabv.NewDHDD(fm.blockCount, time.Now().UnixNano()).
		ParallelInitBuffers(NewRawRsaBlock(fm.rr))

	// read & merge
	fmt.Println("Parallel merging...")
	for i := 0; i < fm.blockCount; i++ {
		// create block
		var block = RawRsaBlock{Data: fm.NextBlockData(), Tag: fm.NextBlockTag(), Key: fm.rr}
		dhdd.ParallelMergeBlock(i, &block)
	}

	fmt.Println("Parallel checking...")
	return dhdd.ParallelCheckAllBuffers()
}
