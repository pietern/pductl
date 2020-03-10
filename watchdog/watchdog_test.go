package watchdog

import (
	"testing"
	"time"
)

func TestWatchdog(t *testing.T) {
	timeout := 10 * time.Millisecond
	w := NewWatchdog(timeout)
	defer w.Stop()

	// Kick triggers presence signal
	w.Kick()
	if <-w.C != Present {
		t.Errorf("Expected present state")
	}

	// Lack of kick triggers absence signal
	t1 := time.Now()
	w.Kick()
	<-w.C
	state := <-w.C
	t2 := time.Now()
	if state != Absent {
		t.Errorf("Expected absent state")
	}

	if t2.Sub(t1) < timeout {
		t.Errorf("Expected delay before absent state")
	}
}
