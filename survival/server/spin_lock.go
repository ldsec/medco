package survivalserver

import "sync/atomic"

type lockState int32

const (
	available lockState = 0
	locked    lockState = 1
)

type Spin struct {
	lock *lockState
}

func NewSpin() *Spin {
	return &Spin{lock: new(lockState)}
}

func (spin *Spin) Lock() {

	for cond := true; cond; {
		if val := atomic.LoadInt32((*int32)(spin.lock)); val == int32(available) {

			cond = !atomic.CompareAndSwapInt32((*int32)(spin.lock), val, int32(locked))
		}

	}
}

func (spin *Spin) Unlock() {

	for cond := true; cond; {
		if val := atomic.LoadInt32((*int32)(spin.lock)); val == int32(available) {
			panic("Unlock an unlocked spin")
		} else {
			cond = !atomic.CompareAndSwapInt32((*int32)(spin.lock), val, int32(available))
		}
	}

}
