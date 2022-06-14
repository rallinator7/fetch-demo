package pubsub

import (
	"github.com/nats-io/nats.go"
)

type Streamer interface {
	Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error)
}

type Publisher struct {
	stream Streamer
}

func NewPublisher(streamer Streamer) Publisher {
	return Publisher{
		stream: streamer,
	}
}

func (publisher *Publisher) Publish(subject string, content []byte) error {
	_, err := publisher.stream.Publish(subject, content)
	if err != nil {
		return err
	}

	return nil
}
