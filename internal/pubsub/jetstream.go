package pubsub

import "github.com/nats-io/nats.go"

func NewJetstream(url string) (nats.JetStreamContext, error) {
	nc, _ := nats.Connect(url)

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js, nil
}
