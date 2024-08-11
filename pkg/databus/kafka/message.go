package kafka

import "time"

type TopicMessage struct {
	Timestamp time.Time
	Topic     string
	Partition int32
	Offset    int64
	Msg       []byte
}
