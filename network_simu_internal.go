package network

func (net *network) validate(id int) {
	if len(net.ends) <= id || id < 0 {
		panic(errEndpointNotExists)
	}
}

func (net *network) getEndpoint(id int) *endpoint {
	net.validate(id)
	return net.ends[id]
}

func (net *network) service(msg message) {
	if err := net.strategies.after(msg.From, msg.To); err != nil {
		msg.Ack <- err
		return
	}

	end := net.getEndpoint(msg.To)

	end.handler.handleMessage(msg.From, msg.Data)
	msg.Ack <- nil
}

func (net *network) call(msg *message) error {
	net.validate(msg.From)
	net.validate(msg.To)

	if err := net.strategies.before(msg.From, msg.To); err != nil {
		return err
	}

	net.getEndpoint(msg.From).increate()
	net.link <- *msg
	return <-msg.Ack
}

func (net *network) isLongDelay() bool {
	net.mutex.RLock()
	defer net.mutex.RUnlock()

	return net.longDelay
}

func (net *network) isReliable() bool {
	net.mutex.RLock()
	defer net.mutex.RUnlock()

	return net.reliable
}
