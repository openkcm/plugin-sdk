package catalog

import (
	"io"
	"time"

	"github.com/zeebo/errs/v2"
	"google.golang.org/grpc"
)

type closerGroup []io.Closer

// Close closes all closers in the group in reverse order.
func (cs closerGroup) Close() error {
	var errs errs.Group
	for i := len(cs) - 1; i >= 0; i-- {
		errs.Add(cs[i].Close())
	}
	return errs.Err()
}

type closerFunc func()

// groupCloserFuncs returns a closerGroup from the given functions.
func groupCloserFuncs(fns ...func()) closerGroup {
	var closers closerGroup
	for _, fn := range fns {
		closers = append(closers, closerFunc(fn))
	}
	return closers
}

// Close implements the io.Closer interface.
func (fn closerFunc) Close() error {
	fn()
	return nil
}

func gracefulStopWithTimeout(s *grpc.Server, d time.Duration) bool {
	done := make(chan struct{})

	go func() {
		s.GracefulStop()
		close(done)
	}()

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-done:
		return true
	case <-t.C:
		s.Stop()
		return false
	}
}
