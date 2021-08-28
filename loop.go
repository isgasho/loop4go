package loop4go

import (
	"sync/atomic"
	"time"
)

type Loop interface {
	Running() bool

	Start() bool

	Stop()
}

type loop struct {
	duration time.Duration
	running  int32
	queue    EventQueue
	callback func(loop Loop)
}

func NewLoop(d time.Duration, queue EventQueue, callback func(loop Loop)) Loop {
	var t = &loop{}
	t.duration = d
	t.running = 0
	t.queue = queue
	t.callback = callback
	return t
}

func (this *loop) Running() bool {
	return atomic.LoadInt32(&this.running) == 1
}

func (this *loop) Start() bool {
	if this.duration <= 0 {
		return false
	}

	if old := atomic.SwapInt32(&this.running, 1); old != 0 {
		return false
	}

	this.enqueue()

	return true
}

func (this *loop) Stop() {
	if old := atomic.SwapInt32(&this.running, 0); old != 1 {
		return
	}
}

func (this *loop) enqueue() {
	if this.Running() {
		after(this.duration, this.queue, this.exec)
	}
}

func (this *loop) exec() {
	if this.Running() {
		defer this.enqueue()
	}
	this.callback(this)
}

func after(d time.Duration, q EventQueue, callback func()) {
	time.AfterFunc(d, func() {
		if q != nil {
			q.Enqueue(callback)
		} else {
			callback()
		}
	})
}
