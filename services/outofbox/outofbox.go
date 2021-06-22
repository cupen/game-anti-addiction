package outofbox

import (
	"fmt"
	"log"
	"sync"

	"github.com/cupen/game-anti-addiction/auth"
	"github.com/cupen/game-anti-addiction/behavior"
	"github.com/cupen/game-anti-addiction/idcard"
	"github.com/cupen/game-anti-addiction/services/backend"
	"github.com/cupen/game-anti-addiction/services/backend/redisstream"
	"github.com/cupen/game-anti-addiction/services/consumer"
)

const (
	types_idcard   = "idcard"
	types_behavior = "behavior"
)

type GameAntiAddiction struct {
	client    *auth.Client
	queues    map[string]backend.Backend
	consumers map[string]*consumer.Runner
	mux       sync.Mutex
}

type QueryCallback func(resp *idcard.QueryResponse)

func New(c *auth.Client, redisUrl string) (*GameAntiAddiction, error) {
	queues := map[string]backend.Backend{}
	for _, name := range []string{types_behavior, types_idcard} {
		backend, err := redisstream.New(redisUrl, name)
		if err != nil {
			return nil, err
		}
		queues[name] = backend
	}
	return &GameAntiAddiction{
		queues: queues,
	}, nil
}

func (gaa *GameAntiAddiction) Start(cb QueryCallback) error {
	gaa.mux.Lock()
	defer gaa.mux.Unlock()
	return gaa.start(cb)
}

func (gaa *GameAntiAddiction) Stop() {
	gaa.mux.Lock()
	defer gaa.mux.Unlock()
	gaa.stopConsumers()
}

func (gaa *GameAntiAddiction) GetBehaviorQueue() backend.Backend {
	return gaa.queues[types_behavior]
}

func (gaa *GameAntiAddiction) GetIDCardQueue() backend.Backend {
	return gaa.queues[types_idcard]
}
func (gaa *GameAntiAddiction) GetClient() *auth.Client {
	return gaa.client
}

func (gaa *GameAntiAddiction) start(cb QueryCallback) error {
	if gaa.hasConsumerRunning() {
		return fmt.Errorf("already started")
	}

	consumers, err := gaa.newConsumers(cb)
	if err != nil {
		return err
	}
	gaa.startConsumers(consumers)
	return nil
}

func (gaa *GameAntiAddiction) hasConsumerRunning() bool {
	return gaa.consumers != nil
}

func (gaa *GameAntiAddiction) newConsumers(cb QueryCallback) (map[string]*consumer.Runner, error) {
	consumerFuncs := map[string]consumer.ConsumerFunc{
		types_behavior: behavior.ConsumerFunc(gaa.client, 128, 100),
		types_idcard:   idcard.ConsumerFunc(gaa.client, 100, cb),
	}

	consumers := map[string]*consumer.Runner{}
	for _, name := range []string{types_behavior, types_idcard} {
		backend := gaa.queues[name]
		consumerFunc := consumerFuncs[name]
		consumers[name] = consumer.New(backend, consumerFunc)
	}
	return consumers, nil
}

func (gaa *GameAntiAddiction) startConsumers(consumers map[string]*consumer.Runner) {
	for name, consumer := range consumers {
		consumer.Start()
		log.Printf("[game-anti-addiction] consumer:%s started, ", name)
	}
	gaa.consumers = consumers
}

func (gaa *GameAntiAddiction) stopConsumers() {
	if gaa.consumers == nil {
		return
	}
	for name, consumer := range gaa.consumers {
		consumer.StopGracefully()
		log.Printf("[game-anti-addiction] consumer:%s stopped, ", name)
	}
	gaa.consumers = nil
}
