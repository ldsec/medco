package survivalserver

import (
	"math"
	"fmt"
	"sync"
	"sync/atomic"
)

const minusOne32 = int32(-1)

//this is an ad-hoc barrier mechanism wrapping a waitgroup
type Barrier struct {
	condition int32
	value     int32
	waitGroup sync.WaitGroup
}

func NewBarrier(condition int) (barrier *Barrier, err error) {
	if condition < 1 {
		err = fmt.Errorf("The condition number must be at least 1, here %d", condition)
		return
	}
	if condition > math.MaxInt32 {
		err = fmt.Errorf("%d exceeds max int32", condition)
		return
	}

	barrier = &Barrier{
		condition: int32(condition),
	}
	return

}

func (barrier *Barrier) Add(delta int32) {
	barrier.waitGroup.Add(int(delta))
	atomic.AddInt32(&(barrier.value), delta)
	//check for overflow ??

}

func (barrier *Barrier) Done() {
	barrier.waitGroup.Done()
	atomic.AddInt32(&(barrier.value), minusOne32)
	//check for overflow or negative value ?
}

func (barrier *Barrier) ConditionalWait() {
	var conditionExceeded bool
	var old int32
	for consistent := false; !consistent; {
		old = atomic.LoadInt32(&(barrier.value))
		conditionExceeded = old >= barrier.condition
		//looks like a CAS without swap, maybe a mutex instead is more appropriate ??
		consistent = atomic.LoadInt32(&(barrier.value)) == old
	}

	if conditionExceeded {
		barrier.waitGroup.Wait()

		for consistent := false; !consistent; {
			old = atomic.LoadInt32(&(barrier.value))
			//this looks more appropriate here, but still a mutex is simpleer to use ?
			consistent = atomic.CompareAndSwapInt32(&(barrier.value), old, int32(0))
		}
	} //don't wait otherwise !!
}

func (barrier *Barrier) AbsoluteWait() {
	barrier.waitGroup.Wait()
	for consistent := false; !consistent; {
		old := atomic.LoadInt32(&(barrier.value))
		//this looks more appropriate here, but still a mutex is simpleer to use ?
		consistent = atomic.CompareAndSwapInt32(&(barrier.value), old, int32(0))
	}
}