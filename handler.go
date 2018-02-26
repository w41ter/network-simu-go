package network

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type handler struct {
	id  int
	net Network

	mutex    sync.Mutex
	callback endCallback
}

func createHandler(id int) *handler {
	return &handler{id: id}
}

func (h *handler) Call(to int, data []byte) error {
	return h.net.Call(h.id, to, data)
}

func (h *handler) ID() int {
	return h.id
}

func (h *handler) BindReceiver(cb endCallback) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.callback = cb
}

func (h *handler) handleMessage(from int, data []byte) {
	callback := h.getCallback()
	if callback != nil {
		callback(from, data)
	} else {
		log.Debugf("%d: ignore message from: %d", h.id, from)
	}
}

// BindNetwork must call for initialize.
func (h *handler) bindNetwork(net Network) {
	h.net = net
}

func (h *handler) GetCount() uint64 {
	return h.net.GetCount(h.id)
}

func (h *handler) Disable() {
	h.net.Disable(h.id)
}

func (h *handler) Enable() {
	h.net.Enable(h.id)
}

func (h *handler) IsEnable() bool {
	return h.net.IsEnable(h.id)
}

func (h *handler) getCallback() endCallback {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.callback
}

func (h *handler) Close() {
	/* ignore */
}
