package network

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func IgnoreCallback(from int, data []byte) {
	/* ignore */
}

func TestBasic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	builder := CreateBuilder()
	call1From := -1
	call1Data := make([]byte, 0)

	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		call1From = from
		call1Data = data

		wg.Done()
	})

	call2From := -1
	builder.AddEndpoint().BindReceiver(func(from int, data []byte) {
		call2From = from

		wg.Done()
	})

	peer3 := builder.AddEndpoint()
	peer3.BindReceiver(IgnoreCallback)

	builder.Build()
	if err := peer3.Call(peer1.ID(), []byte{0x1}); err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	if call1From != peer3.ID() || call2From != -1 {
		t.Fatalf("want call From: %d get: %d and peer 2: %d",
			peer3.ID(), call1From, call2From)
	}

	if len(call1Data) != 1 || call1Data[0] != 0x1 {
		t.Fatalf("want get Data From: %d", peer1.ID())
	}
}

func TestDisconnect(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	builder := CreateBuilder()
	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		fmt.Printf("run callback: %v\n", data)
		if len(data) == 0 {
			panic("peer is reachable")
		}
		wg.Done()
	})

	peer2 := builder.AddEndpoint()
	peer2.BindReceiver(IgnoreCallback)

	net := builder.Build()
	net.Disable(peer2.ID())

	// local not reachable
	err := peer2.Call(peer1.ID(), []byte{0x1})
	if err != errPeerNotReachable {
		t.Fatal("peer is reachable")
	}

	net.Enable(peer2.ID())

	err = peer2.Call(peer1.ID(), []byte{0x2})
	if err != nil {
		t.Fatalf("call it failed")
	}
	wg.Wait()

	net.Disable(peer1.ID())

	// remote not reachable
	err = peer2.Call(peer1.ID(), []byte{})
	if err != errPeerNotReachable {
		t.Fatal("peer is reachable")
	}
	time.Sleep(1 * time.Second)
}

func TestCounts(t *testing.T) {
	runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup
	wg.Add(15)

	builder := CreateBuilder()
	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		wg.Done()
	})

	peer2 := builder.AddEndpoint()
	peer2.BindReceiver(IgnoreCallback)

	net := builder.Build()
	for i := 0; i < 15; i++ {
		go peer2.Call(peer1.ID(), []byte{0x1})
	}

	wg.Wait()
	if n := net.GetCount(peer2.ID()); n != 15 {
		t.Fatalf("wrong GetCount() %v, expected 15\n", n)
	}
}

func TestConcurrentMany(t *testing.T) {
	nclients := 20
	nrpcs := 10

	var wg sync.WaitGroup
	wg.Add(nrpcs * nclients)
	builder := CreateBuilder()

	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		wg.Done()
	})

	peers := []Handler{}
	for i := 0; i < nclients; i++ {
		peer := builder.AddEndpoint()
		peer.BindReceiver(IgnoreCallback)
		peers = append(peers, peer)
	}

	builder.Build()

	for _, peer := range peers {
		go func(peer Handler) {
			for j := 0; j < nrpcs; j++ {
				go peer.Call(peer1.ID(), []byte{0x1})
			}
		}(peer)
	}

	wg.Wait()
}

func TestUnreliable(t *testing.T) {
	nclients := 20
	nrpcs := 10

	var wg sync.WaitGroup
	wg.Add(nrpcs * nclients)

	builder := CreateBuilder()
	var count uint64
	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		atomic.AddUint64(&count, 1)
		wg.Done()
	})

	peers := []Handler{}
	for i := 0; i < nclients; i++ {
		peer := builder.AddEndpoint()
		peers = append(peers, peer)
	}

	net := builder.Build()
	net.SetReliable(false)
	for _, peer := range peers {
		go func(peer Handler) {
			for j := 0; j < nrpcs; j++ {
				err := peer.Call(peer1.ID(), []byte{0x1})
				if err != nil && err == errTimeout {
					wg.Done()
				}
			}
		}(peer)
	}

	wg.Wait()
	if count >= uint64(nrpcs*nclients) {
		t.Fatalf("count great than: %v", nrpcs*nclients)
	}
}

func TestBenchmark(t *testing.T) {
	runtime.GOMAXPROCS(4)

	n := 100000
	var wg sync.WaitGroup
	wg.Add(n)

	builder := CreateBuilder()
	peer1 := builder.AddEndpoint()
	peer1.BindReceiver(func(from int, data []byte) {
		wg.Done()
	})

	peers := make([]Handler, 0)
	for i := 0; i < n; i++ {
		peer := builder.AddEndpoint()
		peers = append(peers, peer)
	}

	builder.Build()
	t0 := time.Now()
	for iters := 0; iters < n; iters++ {
		go func(i int) {
			peers[i].Call(peer1.ID(), []byte{0x1})
		}(iters)
	}
	wg.Wait()
	fmt.Printf("%v for %v\n", time.Since(t0), n)
}

func TestCheck(t *testing.T) {
	builder := CreateAliveBuilder(100, 50)
	var rwg sync.WaitGroup
	var wwg sync.WaitGroup
	rwg.Add(1)
	wwg.Add(1)
	end1 := builder.AddEndpoint(func(end int) {
		rwg.Done()
	}, func(end int) {
		wwg.Done()
	})

	end2 := builder.AddEndpoint(func(end int) {

	}, func(end int) {

	})

	builder.Build()

	rwg.Wait()
	wwg.Wait()

	end1.Close()
	end2.Close()
}
