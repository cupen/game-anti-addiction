package consumer

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/cupen/game-anti-addiction/services"
	"go.uber.org/ratelimit"
)

type consumerFunc func([][]byte) error

type Runner struct {
	backend     services.Backend
	consumer    consumerFunc
	flagRunning int32
}

func New(backend services.Backend, consumer consumerFunc) *Runner {
	if backend == nil || consumer == nil {
		panic(fmt.Errorf("invalid params: backend or consumer must be non-nil"))
	}
	return &Runner{
		backend:  backend,
		consumer: consumer,
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
