package backend

import "time"

type Backend interface {
	Write([]byte) error
	Read(int, time.Duration) ([][]byte, error)
}
