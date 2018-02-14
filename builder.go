package network

type endCallback func(from int, data []byte)

// Builder used by constom to build network
type Builder interface {
	// AddEndpoint add new end to network, and return it id.
	AddEndpoint(cb endCallback) int
	// Build create instance of Network.
	Build() Network
}

type builder struct {
	ends           []*endpoint
	longDelay      bool
	reliable       bool
	longReordering bool
}

// CreateBuilder create instance of Builder.
func CreateBuilder() Builder {
	return &builder{
		ends:           []*endpoint{},
		longDelay:      false,
		reliable:       true,
		longReordering: false,
	}
}

func (b *builder) AddEndpoint(cb endCallback) int {
	newID := len(b.ends)
	end := createEndpoint(newID, cb)
	b.ends = append(b.ends, end)
	return newID
}

func (b *builder) Build() Network {
	return createNetwork(b)
}
