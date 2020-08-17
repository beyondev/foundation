package asio

import (
	"testing"
)

func TestParallel(t *testing.T) {
	const cont = 100000
	p := NewParallel(100, 10)

	errCh := make(chan error, cont*2)
	for i := 0; i < cont; i++ {
		p.Put(func() error {
			return nil
		}, errCh)
	}

	for i := 0; i < cont; i++ {
		p.Put(func() error {
			errCh <- nil
			return nil
		}, nil)
	}

	for i := 0; i < cont; i++ {
		<-errCh
	}

	p.Stop()
}
