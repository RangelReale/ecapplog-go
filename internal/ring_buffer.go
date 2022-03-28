package internal

import (
	"sync"
)

type RingBuffer struct {
	inputChannel  <-chan interface{}
	outputChannel chan interface{}
	lock          *sync.RWMutex
	dropMsgFn     func(m interface{})
}

func NewRingBuffer(inputChannel <-chan interface{}, outputChannel chan interface{}) *RingBuffer {
	return &RingBuffer{inputChannel, outputChannel, &sync.RWMutex{}, nil}
}

func NewRingBufferWithDropFn(inputChannel <-chan interface{}, outputChannel chan interface{}, dropMsgFn func(m interface{})) *RingBuffer {
	return &RingBuffer{inputChannel, outputChannel, &sync.RWMutex{}, dropMsgFn}
}

func (r *RingBuffer) SetOutputChannel(out chan interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.outputChannel = out
}

func (r *RingBuffer) GetOutputChannel() chan interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.outputChannel
}

func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		r.lock.Lock()
		select {
		case r.outputChannel <- v:
		default:
			ov := <-r.outputChannel
			r.outputChannel <- v
			// drop older message
			if r.dropMsgFn != nil {
				r.dropMsgFn(ov)
			}
		}
		r.lock.Unlock()
	}
	close(r.outputChannel)
}
