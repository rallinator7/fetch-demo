package pubsub

import (
	"github.com/nats-io/nats.go"
)

type EventHandler func(event []byte)

type Queuer interface {
	QueueSubscribe(subj string, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
}

type Subscriber struct {
	queuer Queuer
}

type SubscriberOptions struct {
	Options []nats.SubOpt
}

func NewSubscriber(js nats.JetStreamContext) Subscriber {
	return Subscriber{
		queuer: js,
	}
}

func (subscriber *Subscriber) SubscribeToQueue(subject string, queue string, eventer EventHandler) error {
	handler := func(m *nats.Msg) {
		data := m.Data
		eventer(data)
	}

	_, err := subscriber.queuer.QueueSubscribe(subject, queue, handler)
	if err != nil {
		return err
	}

	return nil
}
