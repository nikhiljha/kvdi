package pa

/*
#cgo pkg-config: libpulse
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <pulse/def.h>
#include <pulse/operation.h>
*/
import "C"

import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// newSuccessChan creates a new channel for retrieving the response of
// a function using `pa_context_success_cb_t`. A pointer that can be passed
// to the C callback is returned as well. It must be unrefed after usage.
func newSuccessChan() (chan bool, unsafe.Pointer) {
	// create a channel for the response
	resChan := make(chan bool)

	// save a pointer to the channel
	chPtr := gopointer.Save(resChan)

	return resChan, chPtr
}

// newIndexChan creates a new channel for retrieving the response of a function
// using `pa_context_index_cb_t`. A pointer that can be passed to the C callback
// is returned as well. It must be unrefed after usage.
func newIndexChan() (chan int, unsafe.Pointer) {
	// create a channel for the response
	resChan := make(chan int)

	// save a pointer to the channel
	chPtr := gopointer.Save(resChan)

	return resChan, chPtr
}

// waitForFinish will block until the provided success operation is complete.
// If the operation is cancelled or fails, the provided error is returned.
func waitForFinish(op *C.pa_operation, resChan chan bool, failErr error) error {
	for {
		select {
		// Listen for a response on the channel
		case success := <-resChan:
			if !success {
				return failErr
			}
			return nil
		default:
			switch C.pa_operation_get_state((*C.pa_operation)(op)) {
			case C.PA_OPERATION_RUNNING:
				continue
			case C.PA_OPERATION_CANCELLED:
				return failErr
			case C.PA_OPERATION_DONE: // This would come in on the response channel
			}
		}
	}
}

// waitForIndexFinis is similar to waitForFinish, except it wait for an index
// on the channel and then returns the value or the provided error.
func waitForIndexFinish(op *C.pa_operation, resChan chan int, failErr error) (int, error) {
	for {
		select {
		// Listen for a response on the channel
		case res := <-resChan:
			return res, nil
		default:
			switch C.pa_operation_get_state((*C.pa_operation)(op)) {
			case C.PA_OPERATION_RUNNING:
				continue
			case C.PA_OPERATION_CANCELLED:
				return 0, failErr
			case C.PA_OPERATION_DONE: // This would come in on the response channel
			}
		}
	}
}
