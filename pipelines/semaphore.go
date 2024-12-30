package pipelines

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, maxReq),
	}
}

func (that *Semaphore) Acquire() {
	that.ch <- struct{}{}
}

func (that *Semaphore) Release() {
	<-that.ch
}
