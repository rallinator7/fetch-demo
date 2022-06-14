package payer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/pubsub"
)

type Queuer interface {
	SubscribeToQueue(subject string, queue string, eventer pubsub.EventHandler) error
}

type Adder interface {
	AddPayer(payer internal.Payer) error
}

type Logger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type Subscriber struct {
	queuer Queuer
	adder  Adder
	logger Logger
}

func NewSubscriber(queuer Queuer, adder Adder, logger Logger) Subscriber {
	return Subscriber{
		queuer: queuer,
		adder:  adder,
		logger: logger,
	}
}

func (subscriber *Subscriber) PayerAddedEvent(ctx context.Context, stream string, queue string) func() error {
	return func() error {
		errchan := make(chan error)
		defer close(errchan)

		subscriber.logger.Infow(
			"listening for events",
			"event", "payerAdded",
		)

		payerAdded := func(bytes []byte) {
			var event internal.PayerAdded

			err := json.Unmarshal(bytes, &event)
			if err != nil {
				errchan <- err
			}

			err = subscriber.adder.AddPayer(event.Payer)
			if err != nil {
				subscriber.logger.Errorw(
					"could not add payer",
					"error", err,
				)
				return
			}

			subscriber.logger.Infow(
				"payer added",
				"id", event.Id,
				"name", event.Name,
			)
		}

		subject := fmt.Sprintf("%s.added", stream)

		err := subscriber.queuer.SubscribeToQueue(subject, queue, payerAdded)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		case err := <-errchan:
			return err
		}
	}
}
