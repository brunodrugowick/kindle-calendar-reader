package scheduler

import (
	"time"
)

// Schedule is just another thing I won't be using for now, but I'll keep it here because I want.
// It schedules a function to run every inSeconds seconds approximately.
func Schedule(f func(), inSeconds int64) {
	timeout := inSeconds
	go func() {
		for true {
			if timeout <= 0 {
				f()
				timeout = 0 + inSeconds
			}
			time.Sleep(time.Second)
			timeout -= 1
		}
	}()
}
