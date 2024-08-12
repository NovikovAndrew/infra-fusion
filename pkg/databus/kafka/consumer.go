package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

type ConsumerMessageHandler func(ctx context.Context, payload TopicMessage) error

func (d *Databus) Consume(ctx context.Context, consumerGroupId string, topic string, handler ConsumerMessageHandler) error {
	cg, err := sarama.NewConsumerGroupFromClient(consumerGroupId, d.client)
	if err != nil {
		return err
	}

	defer func() {
		errClose := cg.Close()
		if errClose != nil {
			d.logger.Warn(
				"[Databus] error by close consumer group",
				slog.String("errorClose", errClose.Error()),
			)
		}
	}()

	return cg.Consume(ctx, []string{topic}, &saramaConsumeHandler{messageHandler: handler})
}

type saramaConsumeHandler struct {
	messageHandler ConsumerMessageHandler
}

func (h *saramaConsumeHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *saramaConsumeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *saramaConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return errors.New("[Databus] message channel was closed")
			}

			if err := h.messageHandler(ctx, TopicMessage{
				Timestamp: message.Timestamp,
				Topic:     message.Topic,
				Partition: message.Partition,
				Offset:    message.Offset,
				Msg:       message.Value,
			}); err != nil {
				slog.Error(
					"[Databus] error processing consumer message",
					slog.String("error", err.Error()),
					slog.String("topic", message.Topic),
					slog.Any("message", message.Value),
					slog.Int64("partition", int64(message.Partition)),
					slog.Int64("offset", message.Offset),
				)
				return err
			}

			session.MarkMessage(message, "")
		case <-ctx.Done():
			return nil
		}
	}
}
