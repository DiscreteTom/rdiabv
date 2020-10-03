package main

import (
	"fmt"
	"time"

	"github.com/DiscreteTom/rdiabv"
)

// Define block count
const n = 10000
const timesForHTRM = 10

func main() {
	runDHDD()
	runHTRM()
}

func runDHDD() {
	// Init DHDD
	dhdd := rdiabv.NewDHDD(n, time.Now().UnixNano()).
		InitBuffers(rdiabv.NewDefaultBlock()) // Use default block to init buffers

	// Merge blocks
	for i := 0; i < n; i++ {
		dhdd.MergeBlock(i, rdiabv.DefaultBlockGenerator())
	}

	// Check
	fmt.Println(dhdd.CheckAllBuffers()) // => true
}

func runHTRM() {
	// Init HTRM
	htrm := rdiabv.NewHTRM(timesForHTRM).
		InitBuffers(rdiabv.NewDefaultBlock()) // Use default block to init buffers

	// Merge blocks
	for i := 0; i < n; i++ {
		htrm.MergeBlock(rdiabv.DefaultBlockGenerator())
	}

	// Check
	fmt.Println(htrm.CheckAllBuffers()) // => true
}
