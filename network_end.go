package network

func (net *network) GetCount(id int) uint64 {
	return net.getEndpoint(id).total()
}

func (net *network) Disable(id int) {
	net.getEndpoint(id).disable()
}

func (net *network) Enable(id int) {
	net.getEndpoint(id).enable()
}

func (net *network) IsEnable(id int) bool {
	return net.getEndpoint(id).isEnable()
}
