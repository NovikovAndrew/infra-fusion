package kafka

import (
	"context"
)

func (d *Databus) ProduceOne(ctx context.Context, topic string, message []byte) error {
	panic("implement me")
}

func (d *Databus) ProduceMany(ctx context.Context, topic string, messages [][]byte) error {
	panic("implement me")
}
