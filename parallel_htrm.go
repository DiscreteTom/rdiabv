package rdiabv

import (
	"math/rand"
	"sync"
)

// ParallelInitBuffers will initialize HTRM buffers by coping the given block in parallel.
// Since the times of HTRM won't be too big, the goroutine count equals HTRM times.
func (htrm *HTRM) ParallelInitBuffers(block Block) *HTRM {
	var wg sync.WaitGroup

	htrm.buffers = make([][]Block, htrm.times)
	for i := 0; i < htrm.times; i++ {
		wg.Add(1)
		go htrm.initSingleTimeBuffers(i, block, &wg)
	}
	wg.Wait()
	return htrm
}

// Init buffers of time i.
func (htrm *HTRM) initSingleTimeBuffers(i int, block Block, wg *sync.WaitGroup) {
	defer wg.Done()

	htrm.buffers[i] = make([]Block, ValuePerTimes)
	// init buffers
	for j := 0; j < ValuePerTimes; j++ {
		htrm.buffers[i][j] = block.Copy() // copy value
	}
}

// ParallelMergeBlock will merge the given block to many buffers in parallel.
// Since the times of HTRM won't be too big, the goroutine count equals HTRM times.
func (htrm *HTRM) ParallelMergeBlock(block Block) {
	var wg sync.WaitGroup

	for i := 0; i < htrm.times; i++ {
		wg.Add(1)
		go htrm.mergeSingleBlock(i, block, &wg)
	}
	wg.Wait()
}

// merge single block to a buffer.
func (htrm *HTRM) mergeSingleBlock(i int, block Block, wg *sync.WaitGroup) {
	defer wg.Done()

	// each time, generate a random number from 0 to ValuePerTimes-1
	j := rand.Int() % ValuePerTimes
	htrm.buffers[i][j].Merge(htrm.buffers[i][j], block)
}

// ParallelCheckAllBuffers will check all buffers whether the data field matches the tag field in parallel.
// Since the times of HTRM won't be too big, the goroutine count equals HTRM times.
func (htrm *HTRM) ParallelCheckAllBuffers() bool {
	var wg sync.WaitGroup

	var mismatch = make(chan bool)               // get value when mismatch
	var stoppers = make([]chan bool, htrm.times) // to notify all goroutine to stop
	var done = make(chan bool)                   // to get notification when all goroutine is done
	for i := 0; i < htrm.times; i++ {
		wg.Add(1)
		go htrm.checkSingleTime(i, &wg, mismatch, stoppers[i])
	}
	go waiter(&wg, done)

	select {
	case <-mismatch:
		go notifyAll(stoppers)
		return false
	case <-done:
		return true
	}
}

// check whether a single dimension's buffer is valid.
func (htrm *HTRM) checkSingleTime(i int, wg *sync.WaitGroup, mismatch chan bool, stopper chan bool) {
	defer wg.Done()
	for j := 0; j < ValuePerDimension; j++ {
		select {
		case <-stopper:
			return
		default:
			if htrm.buffers[i][j].Validate() == false {
				mismatch <- true
				return
			}
		}
	}
}
