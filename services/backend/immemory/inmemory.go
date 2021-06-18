package inmemory

import (
	"sync"
	"time"
)

type InMemory struct {
	queue     chan []byte
	queueSize int
	mux       sync.Mutex
}

func New(queueSize int) *InMemory {
	return &InMemory{
		queue:     make(chan []byte, queueSize),
		queueSize: queueSize,
	}
}

func (mem *InMemory) Write(msg []byte) error {
	mem.mux.Lock()
	defer mem.mux.Unlock()
	if len(mem.queue) >= mem.queueSize {
		<-mem.queue
	}
	mem.queue <- msg
	return nil
}

func (mem *InMemory) Read(count int, timeout time.Duration) ([][]byte, error) {
	msgList := [][]byte{}
	isContinue := true
	for isContinue {
		select {
		case msg := <-mem.queue:
			msgList = append(msgList, msg)
			if len(msgList) >= count {
				isContinue = false
			}
		case <-time.After(timeout):
			isContinue = false
		}
	}
	return msgList, nil
}
