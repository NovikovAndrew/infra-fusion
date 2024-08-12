package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
)

func (d *Databus) ProduceOne(ctx context.Context, topic string, message []byte) error {
	producer, err := sarama.NewSyncProducerFromClient(d.client)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := producer.Close(); closeErr != nil {
			d.logger.Warn("[Databus] error by close produce", slog.StringValue(err.Error()))
		}
	}()

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err = producer.SendMessage(producerMessage)
	return err
}
