package network

import (
	"sync/atomic"

	"github.com/thinkermao/network-simu-go/utils/sync"
)

type endpoint struct {
	count   uint64
	enab    *sync.AtomicBool
	handler Handler
}

func createEndpoint(handler Handler) *endpoint {
	return &endpoint{
		count:   0,
		enab:    sync.NewAtomicBool(),
		handler: handler,
	}
}

func (e *endpoint) isEnable() bool {
	return e.enab.IsSet()
}

func (e *endpoint) enable() {
	e.enab.Set()
}

func (e *endpoint) disable() {
	e.enab.UnSet()
}

func (e *endpoint) increate() uint64 {
	return atomic.AddUint64(&e.count, 1)
}

func (e *endpoint) total() uint64 {
	return atomic.LoadUint64(&e.count)
}
