package main

import (
	"fmt"

	"github.com/DiscreteTom/rdiabv"
)

const timesForHTRM = 15

func runHTRM(fm *FileManager) bool {
	fm.StartSession()
	defer fm.EndSession()

	// init HTRM
	fmt.Println("Initializing HTRM...")
	htrm := rdiabv.NewHTRM(timesForHTRM).
		InitBuffers(NewRawRsaBlock(fm.rr))

	// read & merge
	fmt.Println("Merging...")
	for i := 0; i < fm.blockCount; i++ {
		// create block
		var block = RawRsaBlock{Data: fm.NextBlockData(), Tag: fm.NextBlockTag(), Key: fm.rr}
		htrm.MergeBlock(&block)
	}

	fmt.Println("Checking...")
	return htrm.CheckAllBuffers()
}

func parallelRunHTRM(fm *FileManager) bool {
	fm.StartSession()
	defer fm.EndSession()

	// init HTRM
	fmt.Println("Parallel initializing HTRM...")
	htrm := rdiabv.NewHTRM(timesForHTRM).
		ParallelInitBuffers(NewRawRsaBlock(fm.rr))

	// read & merge
	fmt.Println("Parallel merging...")
	for i := 0; i < fm.blockCount; i++ {
		// create block
		var block = RawRsaBlock{Data: fm.NextBlockData(), Tag: fm.NextBlockTag(), Key: fm.rr}
		htrm.ParallelMergeBlock(&block)
	}

	fmt.Println("Parallel checking...")
	return htrm.ParallelCheckAllBuffers()
}
