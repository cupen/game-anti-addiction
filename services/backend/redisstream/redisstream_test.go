package redisstream

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisStream(t *testing.T) {
	url := "redis://127.0.0.1:6379/0"
	obj, err := New(url)
	assert.NoError(t, err)
	assert.NotNil(t, obj)

	t.Cleanup(func() {
		if err := obj.Clear(); err != nil {
			panic(err)
		}
	})

	cases := []int{1, 10, 100, 1000}
	for _, count := range cases {
		name := fmt.Sprintf("batch-%d", count)
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			msg := []byte{1, 2, 3, 4, 5, 6}
			for i := 0; i < count; i++ {
				err := obj.Write(msg)
				assert.NoError(err)
			}

			msgList, err := obj.Read(count, 1*time.Second)
			assert.NoError(err)
			if assert.Equal(count, len(msgList)) {
				for i := 0; i < count; i++ {
					assert.Equal(msg, msgList[i])
				}
			}

			msgList, err = obj.Read(count, 100*time.Millisecond)
			assert.NoError(err)
			assert.Equal(0, len(msgList))
		})
	}
}
