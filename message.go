package network

type message struct {
	From int
	To   int
	Data []byte
	Ack  chan error
}
