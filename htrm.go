package rdiabv

import "math/rand"

// ValuePerTimes defines how many possible values every random-merging has.
const ValuePerTimes = 3

// HTRM is a struct for HTRM algorithm.
type HTRM struct {
	times int

	// Equals the buffer in the paper.
	// buffers[i][j] means the buffer of "the i-th test of random value j".
	buffers [][]Block
}

// NewHTRM will return a new HTRM object with the given times.
func NewHTRM(times int) (htrm *HTRM) {
	htrm = &HTRM{}
	htrm.times = times
	return
}

// InitBuffers will init HTRM buffers by coping the given block.
func (htrm *HTRM) InitBuffers(block Block) *HTRM {
	htrm.buffers = make([][]Block, htrm.times)

	for i := 0; i < htrm.times; i++ {
		htrm.buffers[i] = make([]Block, ValuePerTimes)

		// init buffers
		for j := 0; j < ValuePerTimes; j++ {
			htrm.buffers[i][j] = block.Copy() // copy value
		}
	}
	return htrm
}

// MergeBlock will merge the given block to buffers many times.
func (htrm *HTRM) MergeBlock(index int, block Block) {
	// merge `times` times
	for i := 0; i < htrm.times; i++ {
		// each time, generate a random number from 0 to ValuePerTimes-1
		j := rand.Int() % ValuePerTimes
		htrm.buffers[i][j].Merge(htrm.buffers[i][j], block)
	}
}

// CheckAllBuffers will check all buffers whether the data field matches the tag field
func (htrm *HTRM) CheckAllBuffers() bool {
	for i := 0; i < htrm.times; i++ {
		for j := 0; j < ValuePerTimes; j++ {
			if htrm.buffers[i][j].Validate() == false {
				return false
			}
		}
	}
	return true
}
