package network

import "fmt"

func (n *network) GetCount(id int) uint64 {
	return n.getEndpoint(id).total()
}

func (n *network) Disable(id int) {
	n.getEndpoint(id).disable()
	fmt.Printf("id:%d enable: %v\n", id, n.getEndpoint(id).enab)
}

func (n *network) Enable(id int) {
	n.getEndpoint(id).enable()
	fmt.Printf("id:%d enable: %v\n", id, n.getEndpoint(id).enab)
}

func (n *network) IsEnable(id int) bool {
	return n.getEndpoint(id).isEnable()
}
