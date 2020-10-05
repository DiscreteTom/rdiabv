package rdiabv

import "sync"

// ParallelInitBuffers will initialize DHDD buffers by coping the given block in parallel.
// Since the dimension of DHDD won't be too big, the goroutine count equals DHDD dimension.
func (dhdd *DHDD) ParallelInitBuffers(block Block) *DHDD {
	var wg sync.WaitGroup

	dhdd.buffers = make([][]Block, dhdd.dimension)
	for i := 0; i < dhdd.dimension; i++ {
		wg.Add(1)
		go dhdd.initSingleDimensionBuffers(i, block, &wg)
	}
	wg.Wait()
	return dhdd
}

// Init buffers of dimension i.
func (dhdd *DHDD) initSingleDimensionBuffers(i int, block Block, wg *sync.WaitGroup) {
	defer wg.Done()

	dhdd.buffers[i] = make([]Block, ValuePerDimension)
	// init buffers
	for j := 0; j < ValuePerDimension; j++ {
		dhdd.buffers[i][j] = block.Copy() // copy value
	}
}

// ParallelMergeBlock will merge the given block to many buffers indicated by blockIndexToLogicalPosition in parallel.
// Since the dimension of DHDD won't be too big, the goroutine count equals DHDD dimension.
func (dhdd *DHDD) ParallelMergeBlock(index int, block Block) {
	var wg sync.WaitGroup
	var logicalPosition = dhdd.blockIndexToLogicalPosition[index]
	for i, j := range logicalPosition {
		wg.Add(1)
		go dhdd.mergeSingleBlock(i, j, block, &wg)
	}
	wg.Wait()
}

// merge single block to a buffer.
func (dhdd *DHDD) mergeSingleBlock(i, j int, block Block, wg *sync.WaitGroup) {
	defer wg.Done()
	dhdd.buffers[i][j].Merge(dhdd.buffers[i][j], block)
}

// ParallelCheckAllBuffers will check all buffers whether the data field matches the tag field in parallel.
// Since the dimension of DHDD won't be too big, the goroutine count equals DHDD dimension.
func (dhdd *DHDD) ParallelCheckAllBuffers() bool {
	var wg sync.WaitGroup

	var mismatch = make(chan bool)                   // get value when mismatch
	var stoppers = make([]chan bool, dhdd.dimension) // to notify all goroutine to stop
	var done = make(chan bool)                       // to get notification when all goroutine is done
	for i := 0; i < dhdd.dimension; i++ {
		wg.Add(1)
		go dhdd.checkSingleDimension(i, &wg, mismatch, stoppers[i])
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
func (dhdd *DHDD) checkSingleDimension(i int, wg *sync.WaitGroup, mismatch chan bool, stopper chan bool) {
	defer wg.Done()
	for j := 0; j < ValuePerDimension; j++ {
		select {
		case <-stopper:
			return
		default:
			if dhdd.buffers[i][j].Validate() == false {
				mismatch <- true
				return
			}
		}
	}
}
