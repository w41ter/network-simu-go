package sync

import "sync/atomic"

// NewAtomicBool creates an AtomicBool with default To false
func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}

// NewAtomicBoolWith creates an AtomicBool with given default value
func NewAtomicBoolWith(ok bool) *AtomicBool {
	ab := NewAtomicBool()
	if ok {
		ab.Set()
	}
	return ab
}

// AtomicBool is an atomic Boolean
// Its methods are all atomic, thus safe To be called by
// multiple goroutines simultaneously
// Note: When embedding into a struct, one should always use
// *AtomicBool To avoid copy
type AtomicBool int32

// Set sets the Boolean To true
func (ab *AtomicBool) Set() {
	atomic.StoreInt32((*int32)(ab), 1)
}

// UnSet sets the Boolean To false
func (ab *AtomicBool) UnSet() {
	atomic.StoreInt32((*int32)(ab), 0)
}

// IsSet returns whether the Boolean is true
func (ab *AtomicBool) IsSet() bool {
	return atomic.LoadInt32((*int32)(ab)) == 1
}

// SetTo sets the boolean with given Boolean
func (ab *AtomicBool) SetTo(yes bool) {
	if yes {
		atomic.StoreInt32((*int32)(ab), 1)
	} else {
		atomic.StoreInt32((*int32)(ab), 0)
	}
}

// SetToIf sets the Boolean To new only if the Boolean matches the old
// Returns whether the set was done
func (ab *AtomicBool) SetToIf(old, new bool) (set bool) {
	var o, n int32
	if old {
		o = 1
	}
	if new {
		n = 1
	}
	return atomic.CompareAndSwapInt32((*int32)(ab), o, n)
}

func (ab *AtomicBool) String() string {
	if ab.IsSet() {
		return "true"
	}
	return "false"
}
