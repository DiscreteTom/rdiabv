package main

import (
	"fmt"
	"time"

	"github.com/DiscreteTom/rdiabv"
)

func main() {
	// Define block count
	const n = 10000

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
