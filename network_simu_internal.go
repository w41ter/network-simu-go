package network

import "fmt"

func (n *network) validate(id int) {
	if len(n.ends) <= id || id < 0 {
		panic(errEndpointNotExists)
	}
}

func (n *network) getEndpoint(id int) *endpoint {
	n.validate(id)
	return n.ends[id]
}

func (n *network) service(msg message) {
	if err := n.strategies.after(msg.From, msg.To); err != nil {
		fmt.Printf("strategies: %d => %d err: %v\n", msg.From, msg.To, err)
		return
	}

	end := n.getEndpoint(msg.To)

	end.handler.handleMessage(msg.From, msg.Data)
}

func (n *network) call(msg *message) error {
	n.validate(msg.From)
	n.validate(msg.To)

	if err := n.strategies.before(msg.From, msg.To); err != nil {
		return err
	}

	n.getEndpoint(msg.From).increate()
	n.link <- *msg
	return nil
}

func (n *network) isLongDelay() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.longDelay
}

func (n *network) isReliable() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return n.reliable
}
