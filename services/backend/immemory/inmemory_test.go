package inmemory

import (
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	obj := New(10240)
	msg := []byte{1, 2, 3, 4, 5, 6}
	for i := 0; i < b.N; i++ {
		obj.Write(msg)
	}
}
