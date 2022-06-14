package payer

import (
	"encoding/json"
	"fmt"

	internal "github.com/rallinator7/fetch-demo/internal"
)

type PayerSubject int

const (
	Added PayerSubject = iota
)

type Messager interface {
	Publish(subject string, content []byte) error
}

type Publisher struct {
	Stream   string
	messager Messager
}

func NewPublisher(messager Messager, stream string) Publisher {
	return Publisher{
		Stream:   stream,
		messager: messager,
	}
}

func (publisher *Publisher) PlayerAdded(payer internal.Payer) error {
	added := internal.NewPayerAdded(payer)

	bytes, err := json.Marshal(added)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s", publisher.Stream, Added.String())

	err = publisher.messager.Publish(subject, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (subject PayerSubject) String() string {
	switch subject {
	case Added:
		return "added"
	}

	return "unknown"
}
