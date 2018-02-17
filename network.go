package network

// Network is channel-based tcp net work.
//
// simulates a network that can lose requests, lose replies,
// delay messages, and entirely disconnect particular hosts.
type Network interface {
	Call(from, to int, data []byte) error

	// config
	SetReliable(yes bool)
	SetLongDelays(yes bool)

	// for endpoint
	GetCount(id int) uint64
	Disable(id int)
	Enable(id int)
	IsEnable(id int) bool
}

// Handler is a network handle.
type Handler interface {
	Call(to int, data []byte) error
	BindReceiver(cb endCallback)
	ID() int

	GetCount() uint64
	Disable()
	Enable()
	IsEnable() bool
}
