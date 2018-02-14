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
	if reliable == false {
		// short delay
		ms := rand.Int() % 27
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	if reliable == false && (rand.Int()%1000) < 100 {
		// drop the request, return as if timeout
		return errTimeout
	}

	return nil
}

func (sg *strategies) reachableStrategies(id int) error {
	end := sg.net.getEndpoint(id)
	if !end.isEnable() {
		// simulate no reply and eventual timeout.
		ms := 0
		if sg.net.isLongDelay() {
			// let Raft tests check that leader doesn't send
			// RPCs synchronously.
			ms = rand.Int() % 7000
		} else {
			// many kv tests require the client To try each
			// server in fairly rapid succession.
			ms = rand.Int() % 100
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return errPeerNotReachable
	}
	return nil
}

func (sg *strategies) before(from, to int) error {
	if err := sg.reachableStrategies(from); err != nil {
		return err
	}

	return sg.unreliableStrategies()
}

func (sg *strategies) after(from, to int) error {
	return sg.reachableStrategies(to)
}
