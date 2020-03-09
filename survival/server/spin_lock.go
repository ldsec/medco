package survivalserver

import "sync/atomic"

type lockState int32

//Spin implements a spinlock
type Spin struct {
	lock *lockState
}

// NewSpin spin constructor
func NewSpin() *Spin {
	return &Spin{lock: new(lockState)}
}

// Lock loops until the spin is released, then holds it
func (spin *Spin) Lock() {

	for cond := true; cond; {
		if val := atomic.LoadInt32((*int32)(spin.lock)); val == int32(available) {

			cond = !atomic.CompareAndSwapInt32((*int32)(spin.lock), val, int32(locked))
		}

	}
}

// Unlock releases the spin
func (spin *Spin) Unlock() {

	for cond := true; cond; {
		if val := atomic.LoadInt32((*int32)(spin.lock)); val == int32(available) {
			panic("Unlock an unlocked spin")
		} else {
			cond = !atomic.CompareAndSwapInt32((*int32)(spin.lock), val, int32(available))
		}
	}

}
