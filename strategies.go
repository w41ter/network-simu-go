package network

import (
	"math/rand"
	"time"
)

type strategies struct {
	net *network
}

func createStrategiesHandler(net *network) *strategies {
	return &strategies{
		net: net,
	}
}

func (sg *strategies) unreliableStrategies() error {
	reliable := sg.net.isReliable()
	if !reliable && (rand.Int()%1000) < 100 {
		// drop the request, return as if timeout
		return errTimeout
	}

	return nil
}

func (sg *strategies) longDelayStrategies() {
	longDelay := sg.net.isLongDelay()
	if longDelay {
		ms := 0
		ms = rand.Int() % 7000
		time.Sleep(time.Duration(ms) * time.Millisecond)
	} else {
		// short delay
		ms := rand.Int() % 27
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (sg *strategies) reachableStrategies(id int) error {
	end := sg.net.getEndpoint(id)
	if !end.isEnable() {
		return errPeerNotReachable
	}
	return nil
}

func (sg *strategies) before(from, to int) error {
	sg.longDelayStrategies()
	if err := sg.unreliableStrategies(); err != nil {
		return err
	}

	if err := sg.reachableStrategies(from); err != nil {
		return err
	}

	return sg.reachableStrategies(to)
}

func (sg *strategies) after(from, to int) error {
	return sg.reachableStrategies(to)
}
