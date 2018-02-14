package network

import (
	"sync/atomic"

	"github.com/thinkermao/network-simu-go/utils/sync"
)

type endpoint struct {
	id       int
	count    uint64
	enab     *sync.AtomicBool
	callback endCallback
}

func createEndpoint(id int, callback endCallback) *endpoint {
	return &endpoint{
		id:       id,
		count:    0,
		enab:     sync.NewAtomicBool(),
		callback: callback,
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
