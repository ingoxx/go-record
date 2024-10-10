package cb

import (
	"errors"
	"sync"
	"time"
)

var ErrUnavailable = errors.New("unavailable")

type Group struct {
	wg    sync.WaitGroup
	begin sync.Once
	end   sync.Once
	open  chan bool  // Circuit breaker is open if this channel is closed.
	done  chan bool  // Group is terminated if done is close.
	ch    chan error // Check if the circuit breaker is open, half-open or closed.

	timeout time.Duration
	count   int
	success int
	failure int
}

func New() *Group {
	return &Group{
		ch:      make(chan error),
		done:    make(chan bool),
		open:    make(chan bool),
		timeout: 1 * time.Second,
		failure: 5,
		success: 5,
	}
}

func (g *Group) Stop() {
	g.stop()
}

func (g *Group) Do(fn func() error) error {
	select {
	case <-g.open:
		return ErrUnavailable
	default:
		err := fn()

		g.init() // Lazily inits the background job.

		select {
		case <-g.open:
			return ErrUnavailable
		case g.ch <- err:
		}

		return err
	}
}

func (g *Group) stop() {
	g.end.Do(func() {
		close(g.done)
		g.wg.Wait()
	})
}

func (g *Group) init() {
	g.begin.Do(func() {
		g.wg.Add(1)

		go func() {
			defer g.wg.Done()

			g.worker()
		}()
	})
}

func (g *Group) worker() {
	t := time.NewTimer(0)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// The service recovers after timeout.
			g.open = make(chan bool) // Open to Half-Open.
			g.count = 0
		case <-g.done:
			return
		case err := <-g.ch:
			// Double-check to see if the circuit breaker is opened.
			// The service is healthy.
			if err == nil {
				continue
			}

			// The service is unhealthy.
			// After a certain threshold, circuit breaker becomes Open.
			g.count++
			if g.count > g.failure {
				close(g.open) // Closed to Open.
				t.Reset(g.timeout)
			}
		}
	}
}
