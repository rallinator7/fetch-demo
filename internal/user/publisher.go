package user

import (
	"encoding/json"
	"fmt"

	internal "github.com/rallinator7/fetch-demo/internal"
)

type UserSubject int

const (
	Added UserSubject = iota
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

func (publisher *Publisher) UserAdded(user internal.User) error {
	added := internal.NewUserAdded(user)

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

func (subject UserSubject) String() string {
	switch subject {
	case Added:
		return "added"
	}

	return "unknown"
}
