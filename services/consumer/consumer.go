package consumer

import (
	"log"
	"time"

	"github.com/cupen/game-anti-addiction/services"
)

type Handler func([][]byte) error

func Run(queue services.Backend, h Handler) error {
	isContinue := true
	for i := 0; isContinue; i++ {
		time.Sleep(1 * time.Second)
		msgList, err := queue.Read(128, 1*time.Second)
		if err != nil {
			continue
		}

		if len(msgList) <= 0 {
			continue
		}
		if err := h(msgList); err != nil {
			log.Printf("handle msglist failed: %v", err)
			continue
		}
	}
	return nil
}
