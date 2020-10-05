package rdiabv

import "sync"

// notifyAll will notify all goroutine to stop.
func notifyAll(stoppers []chan bool) {
	for i := 0; i < len(stoppers); i++ {
		stoppers[i] <- true
	}
}

// waiter will wait until all work is done.
func waiter(wg *sync.WaitGroup, done chan bool) {
	wg.Wait()
	done <- true
}
