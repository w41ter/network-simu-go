package network

type endCallback func(from int, data []byte)

// Builder used by constom to build network
type Builder interface {
	// AddEndpoint add new end to network, and return it id.
	AddEndpoint() Handler
	// Build create instance of Network.
	Build() Network
}

type builder struct {
	ends      []*endpoint
	longDelay bool
	reliable  bool
}

// CreateBuilder create instance of Builder.
func CreateBuilder() Builder {
	return &builder{
		ends:      []*endpoint{},
		longDelay: false,
		reliable:  true,
	}
}

func (b *builder) AddEndpoint() Handler {
	newID := len(b.ends)
	h := createHandler(newID)
	end := createEndpoint(h)
	b.ends = append(b.ends, end)

	return h
}

func (b *builder) Build() Network {
	net := createNetwork(b)
	for i := 0; i < len(b.ends); i++ {
		end := b.ends[i]
		end.handler.BindNetwork(net)
	}
	return net
}
