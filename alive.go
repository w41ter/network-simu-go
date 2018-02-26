package network

import (
	"runtime"
	"sync/atomic"
	"time"
)

type writeTimeoutCallback func(end int)
type readTimeoutCallback func(end int)

type remote struct {
	id             int
	write          chan []byte
	read           chan []byte
	readTimeout    int
	writeTimeout   int
	readTimeoutCh  <-chan time.Time
	writeTimeoutCh <-chan time.Time
	readCb         readTimeoutCallback
	writeCb        writeTimeoutCallback
	exit           int32
}

func createRemote(id, readTimeout, writeTimeout int,
	read readTimeoutCallback, write writeTimeoutCallback) *remote {
	r := &remote{
		id:             id,
		read:           make(chan []byte, 1),
		write:          make(chan []byte, 1),
		readCb:         read,
		writeCb:        write,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		readTimeoutCh:  time.After(time.Duration(readTimeout) * time.Millisecond),
		writeTimeoutCh: time.After(time.Duration(writeTimeout) * time.Millisecond),
		exit:           0,
	}

	r.service()

	return r
}

func (r *remote) service() {
	go func() {
		for atomic.LoadInt32(&r.exit) == 0 {
			runtime.Gosched()

			select {
			case <-r.read:
				r.readTimeoutCh = time.After(time.Duration(r.readTimeout) * time.Millisecond)
			case <-r.write:
				r.writeTimeoutCh = time.After(time.Duration(r.writeTimeout) * time.Millisecond)
			default:
				/* ignore */
			}

			if r.readTimeoutCh != nil {
				select {
				case <-r.readTimeoutCh:
					r.readCb(r.id)
					r.readTimeoutCh = nil
				default:
					/* ignore */
				}
			}

			if r.writeTimeoutCh != nil {
				select {
				case <-r.writeTimeoutCh:
					r.writeCb(r.id)
					r.writeTimeoutCh = nil
				default:
					/* ignore */
				}
			}
		}
	}()
}

func (r *remote) handleMessage(data []byte) {
	r.read <- data
}

func (r *remote) sendMessage(data []byte) {
	r.write <- data
}

func (r *remote) close() {
	atomic.SwapInt32(&r.exit, 1)
	close(r.read)
	close(r.write)
}

type aliveHandler struct {
	*handler
	readCb       readTimeoutCallback
	writeCb      writeTimeoutCallback
	readTimeout  int
	writeTimeout int
	remotes      []*remote
}

func createAliveHandler(id, readTimeout, writeTimeout int,
	readCb readTimeoutCallback, writeCb writeTimeoutCallback) *aliveHandler {
	h := &aliveHandler{
		handler:      createHandler(id),
		readCb:       readCb,
		writeCb:      writeCb,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
	return h
}

func (h *aliveHandler) setEndpoints(ends []int) {
	h.remotes = make([]*remote, 0)
	for i := 0; i < len(ends); i++ {
		if ends[i] == h.id {
			continue
		}
		remote := createRemote(ends[i], h.readTimeout, h.writeTimeout,
			h.handleReadTimeout, h.handleWriteTimeout)
		h.remotes = append(h.remotes, remote)
	}
}

func (h *aliveHandler) handleWriteTimeout(end int) {
	if h.writeCb != nil {
		h.writeCb(end)
	}
}

func (h *aliveHandler) handleReadTimeout(end int) {
	if h.readCb != nil {
		h.readCb(end)
	}
}

func (h *aliveHandler) getEndpoint(from int) *remote {
	for i := 0; i < len(h.remotes); i++ {
		if h.remotes[i].id == from {
			return h.remotes[i]
		}
	}
	panic("wrong endpoint id")
}

func (h *aliveHandler) Call(to int, data []byte) error {
	remote := h.getEndpoint(to)
	remote.sendMessage(data)
	return h.handler.Call(to, data)
}

func (h *aliveHandler) handleMessage(from int, data []byte) {
	remote := h.getEndpoint(from)
	remote.handleMessage(data)
	h.handler.handleMessage(from, data)
}

func (h *aliveHandler) Close() {
	for i := 0; i < len(h.remotes); i++ {
		h.remotes[i].close()
	}
}
