package watchdog

import (
	"time"
)

type State int

const (
	Present State = iota
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

	for {
		select {
		case _, ok := <-w.kick:
			if !ok {
				return
			}

			// Reset timer.
			if !t.Stop() {
				<-t.C
			}
			t.Reset(w.timeout)

			// Notify downstream of presence.
			ch <- Present

		case _, ok := <-t.C:
			if !ok {
				return
			}

			// Notify downstream of absence.
			ch <- Absent
		}
	}
}
