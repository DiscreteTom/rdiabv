package main

import (
	"fmt"

	"github.com/DiscreteTom/rdiabv"
)

const timesForHTRM = 10

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
