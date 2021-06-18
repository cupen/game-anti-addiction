package redisstream

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisStream struct {
	redis     *redis.Client
	key       string
	writeOpts struct {
		maxLen int64
	}
	readOpts struct {
		lastMsgID string
	}
}

func New(url string) (*RedisStream, error) {
	key := "game-anti-addiction/queue/v1"
	maxLen := int64(10240)
	return NewWithArgs(url, key, maxLen)
}

func NewWithArgs(urlOrClient interface{}, key string, maxLen int64) (*RedisStream, error) {
	var client *redis.Client
	switch o := urlOrClient.(type) {
	case string:
		opts, err := redis.ParseURL(o)
		if err != nil {
			return nil, err
		}
		client = redis.NewClient(opts)
	case *redis.Client:
		client = o
	default:
		panic(fmt.Errorf("invalid args:urlOrClient, it must be url string or *redis.Client"))
	}

	if key == "" {
		return nil, fmt.Errorf("empty key")
	}
	if maxLen <= 128 {
		return nil, fmt.Errorf("maxLen must be greater than 128")
	}
	rdss := RedisStream{
		redis: client,
		key:   key,
	}
	rdss.setupDefault()
	rdss.writeOpts.maxLen = maxLen
	return &rdss, nil
}

func (r *RedisStream) setupDefault() {
	// write
	r.writeOpts.maxLen = 10240
}

func (r *RedisStream) Clear() error {
	ctx := context.Background()
	keyOfLastId := r._buildKeyOfLastMsgId()
	err := r.redis.Del(ctx, r.key, keyOfLastId).Err()
	if err == redis.Nil {
		return nil
	}
	return nil
}

func (r *RedisStream) Write(msg []byte) error {
	ctx := context.Background()
	args := redis.XAddArgs{
		Stream:       r.key,
		MaxLen:       r.writeOpts.maxLen,
		MaxLenApprox: makeApprox(r.writeOpts.maxLen),
		Values:       []interface{}{"body", msg},
	}
	rs := r.redis.XAdd(ctx, &args)
	// log.Printf("[xwrite] %v", rs.Args())
	return rs.Err()
}

func (r *RedisStream) _buildKeyOfLastMsgId() string {
	return r.key + ":lastMsgId"
}

func (r *RedisStream) getLastMsgId() (string, error) {
	if r.readOpts.lastMsgID != "" {
		return r.readOpts.lastMsgID, nil
	}

	key := r._buildKeyOfLastMsgId()
	ctx := context.Background()
	rs, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return rs, nil
}

func (r *RedisStream) setLastMsgId(id string) error {
	r.readOpts.lastMsgID = id

	key := r._buildKeyOfLastMsgId()
	ctx := context.Background()
	_, err := r.redis.Set(ctx, key, id, 0).Result()
	return err
}

func (r *RedisStream) Read(count int, timeout time.Duration) ([][]byte, error) {
	lastMsgId, err := r.getLastMsgId()
	if err != nil {
		return nil, err
	}
	if lastMsgId == "" {
		lastMsgId = "0-0"
		r.setLastMsgId(lastMsgId)
	}

	ctx := context.Background()
	args := redis.XReadArgs{
		Streams: []string{r.key, lastMsgId},
		Count:   int64(count),
		Block:   timeout,
	}
	rs := r.redis.XRead(ctx, &args)
	// log.Printf("[xread] %v", rs.Args())
	streams, err := rs.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if len(streams) != 1 {
		return nil, fmt.Errorf("must be one stream")
	}

	toBytes := func(v interface{}) []byte {
		switch _v := v.(type) {
		case []byte:
			log.Printf("received bytes :%v", _v)
			return _v
		case string:
			// log.Printf("received string :%v", _v)
			return []byte(_v)
		default:
			return nil
		}
	}

	msgList := [][]byte{}
	for _, msg := range streams[0].Messages {
		for k, v := range msg.Values {
			_v := toBytes(v)
			if _v == nil {
				log.Printf("invalid message k = %s v = %v", k, reflect.TypeOf(v))
				continue
			}
			msgList = append(msgList, _v)
		}
		lastMsgId = msg.ID
	}
	r.setLastMsgId(lastMsgId)
	return msgList, nil
}

func makeApprox(maxLen int64) int64 {
	if maxLen >= 10000 {
		return maxLen / 10 // 10%
	}
	if maxLen > 1000 {
		return (maxLen / 10) * 3 // 30%
	}
	return 256
}
