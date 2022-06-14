package user

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
	AddUser(payer internal.User) error
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

func (subscriber *Subscriber) UserAddedEvent(ctx context.Context, stream string, queue string) func() error {
	return func() error {
		errchan := make(chan error)
		defer close(errchan)

		subscriber.logger.Infow(
			"listening for events",
			"event", "userAdded",
		)

		userAdded := func(bytes []byte) {
			var event internal.UserAdded

			err := json.Unmarshal(bytes, &event)
			if err != nil {
				errchan <- err
			}

			err = subscriber.adder.AddUser(event.User)
			if err != nil {
				subscriber.logger.Errorw(
					"could not add user",
					"error", err,
				)
				return
			}

			subscriber.logger.Infow(
				"user added",
				"id", event.Id,
				"first_name", event.FirstName,
				"last_name", event.LastName,
			)
		}

		subject := fmt.Sprintf("%s.added", stream)

		err := subscriber.queuer.SubscribeToQueue(subject, queue, userAdded)
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
