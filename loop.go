package goey

import (
	"errors"
	"runtime"
	"sync/atomic"
)

var (
	// ErrQuit indicates that the event loop should terminate.  This return
	// will only be used on platforms that expose their loop iteration function
	// in addition to Run.
	ErrQuit = errors.New("quit message")

	// ErrNotRunning indicates that the main loop is not running.
	ErrNotRunning = errors.New("main loop is not running")

	// ErrAlreadyRunning indicates that the main loop is not running.
	ErrAlreadyRunning = errors.New("main loop is already running")
)

var (
	isRunning uint32
)

// Run locks the OS thread to act as a GUI thread, and then iterates over the
// event loop until there are no more instances of MainWindow open.
// If the main loop is already running, this function will return an error.
//
// Modification of the GUI should happen only on the GUI thread.  Any
// modifications should be schedule using a callback on the function Do.
func Run() error {
	// Pin the GUI message loop to a single thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Want to gate entry into the GUI loop so that only one thread may enter
	// at a time.  Since this is supposed to be non-blocking, we can't use
	// a sync.Mutex without a TryLock method.
	if !atomic.CompareAndSwapUint32(&isRunning, 0, 1) {
		return ErrAlreadyRunning
	}
	defer func() {
		atomic.StoreUint32(&isRunning, 0)
	}()

	// Defer to platform-specific code.
	return run()
}

// Do runs the passed function on the GUI thread.  If the event loop is not
// running, this function will return an error.  Any error from the callback will
// also be returned.
//
// Note, this function contains a race-condition, in that the the action may be
// scheduled while the event loop is being terminated, in which case the
// scheduled action may never be run.  Presumably, those actions don't need to
// be run on the GUI thread, so they can be schedule using a different
// mechanism.
func Do(action func() error) error {
	// Check if the event loop is current running.
	if atomic.LoadUint32(&isRunning) == 0 {
		return ErrNotRunning
	}

	// Race-condition here!  Event loop may terminate between previous check
	// and following call, which will block.

	// Defer to platform-specific code.
	return do(action)
}
