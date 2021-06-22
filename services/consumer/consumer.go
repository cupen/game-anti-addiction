package consumer

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"go.uber.org/ratelimit"
)

type Backend interface {
	Write([]byte) error
	Read(int, time.Duration) ([][]byte, error)
}

type ConsumerFunc func([][]byte) error

type Runner struct {
	backend     Backend
	consumer    ConsumerFunc
	flagRunning int32
}

func New(backend Backend, cbfunc ConsumerFunc) *Runner {
	if backend == nil || cbfunc == nil {
		panic(fmt.Errorf("invalid params: backend or consumer must be non-nil"))
	}
	return &Runner{
		backend:  backend,
		consumer: cbfunc,
	}
}

func (c *Runner) Start() {
	go c.Run()
}

func (c *Runner) StopGracefully() {
	c.setRunning(false)
}

func (c *Runner) Run() error {
	limiter := ratelimit.New(100)
	c.setRunning(true)
	for i := 0; c.isRunning(); i++ {
		_ = limiter.Take()
		msgList, err := c.backend.Read(128*10, 1*time.Second)
		if err != nil {
			continue
		}

		if len(msgList) <= 0 {
			continue
		}
		if err := c.consumer(msgList); err != nil {
			log.Printf("handle msglist failed: %v", err)
			continue
		}
	}
	return nil
}

func (c *Runner) setRunning(v bool) {
	if v {
		atomic.StoreInt32(&c.flagRunning, 1)
	} else {
		atomic.StoreInt32(&c.flagRunning, 0)
	}
}

func (c *Runner) isRunning() bool {
	return atomic.LoadInt32(&c.flagRunning) == 1
}
