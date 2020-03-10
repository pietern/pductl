package watchdog

import (
	"time"
)

type State int

const (
	Invalid State = iota
	Present
	Absent
)

type Watchdog struct {
	C <-chan State

	done    chan struct{}
	kick    chan struct{}
	timeout time.Duration
}

func NewWatchdog(timeout time.Duration) *Watchdog {
	ch := make(chan State, 1)
	w := Watchdog{
		C: ch,

		done:    make(chan struct{}),
		kick:    make(chan struct{}),
		timeout: timeout,
	}

	go w.run(ch)

	return &w
}

// Stop stops the watchdog.
//
// The state channel has been closed when this function returns.
//
func (w *Watchdog) Stop() {
	close(w.kick)
	<-w.done
}

// Kick lets the watchdog know its subject is alive.
func (w *Watchdog) Kick() {
	w.kick <- struct{}{}
}

func (w *Watchdog) run(ch chan<- State) {
	defer close(w.done)

	t := time.NewTimer(w.timeout)
	defer t.Stop()

	prev := Invalid
	for {
		select {
		case _, ok := <-w.kick:
			if !ok {
				return
			}

			// Reset timer.
			if !t.Stop() {
				// Try to drain the timer channel in case it has
				// fired but has not yet been processed.
				//
				// Also see https://github.com/golang/go/issues/27169
				//
				select {
				case <-t.C:
				default:
				}
			}
			t.Reset(w.timeout)

			// Notify downstream of presence.
			if prev != Present {
				prev = Present
				ch <- Present
			}

		case _, ok := <-t.C:
			if !ok {
				return
			}

			// Notify downstream of absence.
			if prev != Absent {
				prev = Absent
				ch <- Absent
			}
		}
	}
}
