package goq

import "sync/atomic"

type (
	// Manager interface
	GoqManager interface {
		// Wait for available slot
		Wait()

		// Tag goroutine as complete
		Done()

		// Close manager
		Close()

		// Wait for all goroutines to finish
		WaitAllDone()

		// Returns running count of goroutines
		RunningCount() int32
	}

	goqManager struct {
		// Number of maximum concurrent goroutines
		max int

		// Channel to co-ordinate number of concurrent goroutines
		managerCh chan interface{}

		// Channel to notify of completion of a single goroutine
		doneCh chan bool

		// Channel to notify that all goroutines have completed
		allDoneCh chan bool

		// Allows us to check if its ok to close the manager
		closed bool

		// Contains the number of running goroutines
		runningCount int32
	}
)

// New goqManager
func New(maxGoRoutines int) *goqManager {
	// Initiate the manager object
	c := goqManager{
		max:       maxGoRoutines,
		managerCh: make(chan interface{}, maxGoRoutines),
		doneCh:    make(chan bool),
		allDoneCh: make(chan bool),
	}

	// Fill the manager channel by placeholder values
	for i := 0; i < c.max; i++ {
		c.managerCh <- nil
	}

	// Start the controller to collect all the jobs
	go c.controller()

	return &c
}

// Create the controller to collect all the jobs.
// When a goroutine is finished, we can release a slot for another goroutine.
func (c *goqManager) controller() {
	for {
		// This will block until a goroutine is finished
		<-c.doneCh

		// Say that another goroutine can now start
		c.managerCh <- nil

		// When the closed flag is set,
		// we need to close the manager if it doesn't have any running goroutine
		if c.closed && c.runningCount == 0 {
			break
		}
	}

	// Say that all goroutines are finished, we can close the manager
	c.allDoneCh <- true
}

// Wait until a slot is available for the new goroutine.
// A goroutine have to start after this function.
func (c *goqManager) Wait() {

	// Try to receive from the manager channel. When we have something,
	// it means a slot is available and we can start a new goroutine.
	// Otherwise, it will block until a slot is available.
	<-c.managerCh

	// Increase the running count to help we know how many goroutines are running.
	atomic.AddInt32(&c.runningCount, 1)
}

// Mark a goroutine as finished
func (c *goqManager) Done() {
	// Decrease the number of running count
	atomic.AddInt32(&c.runningCount, -1)
	c.doneCh <- true
}

// Close the manager manually
func (c *goqManager) Close() {
	c.closed = true
}

// Wait for all goroutines are done
func (c *goqManager) WaitAllDone() {
	// Close the manager automatic
	c.Close()

	// This will block until allDoneCh was marked
	<-c.allDoneCh
}

// Returns the number of goroutines which are running
func (c *goqManager) RunningCount() int32 {
	return c.runningCount
}