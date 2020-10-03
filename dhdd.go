package rdiabv

import (
	"math/rand"
)

// DHDD is a struct for DHDD algorithm.
type DHDD struct {
	dimension                   int           // Equals the x in the paper.
	logicalBlockNum             int           // Equals the N in the paper.
	blockIndexToLogicalPosition map[int][]int // Equals the pi in the paper.

	// Equals the buffer in the paper.
	// buffers[i][j] means the buffer of "the value of dimension i is j".
	buffers [][]Block
}

// NewDHDD calculate dimension & logicalBlockNum according to the count of blocks n
// and generate blockIndexToLogicalPosition randomly by using rand.Seed(seed),
// then return the DHDD object. If you are not sure what seed to use, use time.Now().UnixNano().
func NewDHDD(n int, seed int64) (dhdd *DHDD) {
	dhdd = &DHDD{}
	dhdd.getDimensionAndLogicalBlockNum(n)
	dhdd.generateBlockIndexToLogicalPosition(seed)
	return
}

func (dhdd *DHDD) getDimensionAndLogicalBlockNum(n int) {
	dhdd.dimension = 0
	dhdd.logicalBlockNum = 1
	for {
		if dhdd.logicalBlockNum < n {
			dhdd.logicalBlockNum *= 3
			dhdd.dimension++
		} else {
			break
		}
	}
}

func (dhdd *DHDD) generateBlockIndexToLogicalPosition(seed int64) {
	// generate logical positions, each logical position is an int list of length x
	var logicalPositions = cartesianProduct(dhdd.dimension)
	// shuffle logical positions
	rand.Seed(seed)
	rand.Shuffle(dhdd.logicalBlockNum, func(i, j int) {
		logicalPositions[i], logicalPositions[j] = logicalPositions[j], logicalPositions[i]
	})
	// construct result
	dhdd.blockIndexToLogicalPosition = make(map[int][]int)
	for i := 0; i < dhdd.logicalBlockNum; i++ {
		dhdd.blockIndexToLogicalPosition[i] = logicalPositions[i]
	}
}

// InitBuffers will init DHDD buffers by coping the given block.
func (dhdd *DHDD) InitBuffers(block Block) *DHDD {
	dhdd.buffers = make([][]Block, dhdd.dimension)
	for i := 0; i < dhdd.dimension; i++ {
		dhdd.buffers[i] = make([]Block, 3)

		// init buffers
		dhdd.buffers[i][0] = block.Copy() // copy value
		dhdd.buffers[i][1] = block.Copy()
		dhdd.buffers[i][2] = block.Copy()
	}
	return dhdd
}

// MergeBlock will merge the given block to many buffers indicated by blockIndexToLogicalPosition.
func (dhdd *DHDD) MergeBlock(index int, block Block) {
	var logicalPosition = dhdd.blockIndexToLogicalPosition[index]
	for i, j := range logicalPosition {
		dhdd.buffers[i][j].Merge(dhdd.buffers[i][j], block)
	}
}

// CheckAllBuffers will check all buffers whether the data field matches the tag field
func (dhdd *DHDD) CheckAllBuffers() bool {
	for i := 0; i < len(dhdd.buffers); i++ {
		for j := 0; j < len(dhdd.buffers[i]); j++ {
			if dhdd.buffers[i][j].Validate() == false {
				return false
			}
		}
	}
	return true
}
