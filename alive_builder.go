package network

// AliveBuilder used by custom to build network
type AliveBuilder interface {
	// AddEndpoint add new end to network, and return it id.
	AddEndpoint(readCb readTimeoutCallback,
		writeCb writeTimeoutCallback) Handler
	// Build create instance of Network.
	Build() Network
}

type aliveBuilder struct {
	ends                      []*endpoint
	handlers                  []*aliveHandler
	longDelay                 bool
	reliable                  bool
	readTimeout, writeTimeout int
}

// CreateAliveBuilder create instance of AliveBuilder.
func CreateAliveBuilder(readTimeout, writeTimeout int) AliveBuilder {
	return &aliveBuilder{
		ends:         []*endpoint{},
		handlers:     []*aliveHandler{},
		longDelay:    false,
		reliable:     true,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (b *aliveBuilder) AddEndpoint(readCb readTimeoutCallback,
	writeCb writeTimeoutCallback) Handler {
	newID := len(b.ends)
	h := createAliveHandler(newID, b.readTimeout, b.writeTimeout, readCb, writeCb)
	end := createEndpoint(h)
	b.ends = append(b.ends, end)
	b.handlers = append(b.handlers, h)

	return h
}

func (b *aliveBuilder) Build() Network {
	net := createAliveNetwork(b)
	endpoints := net.Endpoints()
	for i := 0; i < len(b.ends); i++ {
		end := b.ends[i]
		end.handler.bindNetwork(net)
		b.handlers[i].setEndpoints(endpoints)
		net.registerEnableListener(b.handlers[i].enableCallback)
	}
	return net
}
