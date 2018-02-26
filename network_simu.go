package network

import (
	"errors"
	"sync"
)

var (
	errPeerNotReachable  = errors.New("peer not reachable")
	errTimeout           = errors.New("send timeout")
	errEndpointNotExists = errors.New("endpoint not exists")
)

type network struct {
	mutex      sync.RWMutex
	longDelay  bool
	reliable   bool
	ends       []*endpoint
	link       chan message
	strategies *strategies
}

func createNetwork(b *builder) *network {
	net := &network{
		longDelay: b.longDelay,
		reliable:  b.reliable,
		ends:      b.ends,
		link:      make(chan message, 1024),
	}

	net.strategies = createStrategiesHandler(net)

	for i := 0; i < len(net.ends); i++ {
		net.ends[i].enable()
	}

	go func() {
		for msg := range net.link {
			/* use copy rather than ref */
			go net.service(msg)
		}
	}()

	return net
}

func createAliveNetwork(b *aliveBuilder) *network {
	net := &network{
		longDelay: b.longDelay,
		reliable:  b.reliable,
		ends:      b.ends,
		link:      make(chan message, 1024),
	}

	net.strategies = createStrategiesHandler(net)

	for i := 0; i < len(net.ends); i++ {
		net.ends[i].enable()
	}

	go func() {
		for msg := range net.link {
			/* use copy rather than ref */
			go net.service(msg)
		}
	}()

	return net
}

func (net *network) Call(from, to int, data []byte) error {
	buf := make([]byte, len(data))
	copy(buf, data)
	msg := message{
		From: from,
		To:   to,
		Data: buf,
	}
	return net.call(&msg)
}

func (net *network) Endpoints() []int {
	var ends []int
	for i := 0; i < len(net.ends); i++ {
		ends = append(ends, net.ends[i].handler.ID())
	}
	return ends
}

func (net *network) SetReliable(yes bool) {
	net.mutex.RLock()
	defer net.mutex.RUnlock()

	net.reliable = yes
}

func (net *network) SetLongDelays(yes bool) {
	net.mutex.RLock()
	defer net.mutex.RUnlock()

	net.longDelay = yes
}
