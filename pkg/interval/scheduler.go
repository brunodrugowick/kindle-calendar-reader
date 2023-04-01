package interval

import (
	"log"
	"reflect"
	"runtime"
	"time"
)

// RunAtInterval runs the given function `f` at the specified interval `d`.
// It stops the execution when the stop channel is closed.
func RunAtInterval(f func(), d time.Duration) chan struct{} {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				log.Printf("Stopping function %s\n", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
				return
			case <-ticker.C:
				f()
			}
		}
	}()
	return stop
}
