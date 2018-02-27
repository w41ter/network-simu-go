package network

func (n *network) GetCount(id int) uint64 {
	return n.getEndpoint(id).total()
}

func (n *network) Disable(id int) {
	n.getEndpoint(id).disable()
}

func (n *network) Enable(id int) {
	for i := 0; i < len(n.enableCallbacks); i++ {
		n.enableCallbacks[i](id)
	}
	n.getEndpoint(id).enable()
}

func (n *network) IsEnable(id int) bool {
	return n.getEndpoint(id).isEnable()
}
