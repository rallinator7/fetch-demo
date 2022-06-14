package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rallinator7/fetch-demo/internal"
)

const (
	streamName     = "PAYER"
	streamSubjects = "PAYER.added"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)

	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = createStream(js)
	if err != nil {
		log.Fatalf(err.Error())
	}

	go publish(js, streamSubjects)

	subscribe(nc, streamSubjects)
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	subjectName = "PAYER.added"
)

func publish(js nats.JetStreamContext, subject string) {
	names := []string{
		"DANNON",
		"UNILEVER",
		"MILLER COORS",
	}

	for {
		for _, name := range names {
			err := addPayer(js, subject, name)
			if err != nil {
				log.Println(err.Error())
			}

			time.Sleep(time.Second)
		}
	}
}

func subscribe(conn *nats.Conn, subject string) {
	// Use a WaitGroup to wait for 10 messages to arrive
	wg := sync.WaitGroup{}
	wg.Add(10)

	// Create a queue subscription on "updates" with queue name "workers"
	_, err := conn.QueueSubscribe(subject, "workers", func(m *nats.Msg) {
		fmt.Printf("%s %d\n", string(m.Data), 1)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a queue subscription on "updates" with queue name "workers"
	_, err = conn.QueueSubscribe(subject, "workers", func(m *nats.Msg) {
		fmt.Printf("%s %d\n", string(m.Data), 2)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Wait for messages to come in
	wg.Wait()
}

func addPayer(js nats.JetStreamContext, subject string, name string) error {
	payer := internal.NewPayer(name)

	payerJSON, _ := json.Marshal(payer)

	_, err := js.Publish(subject, payerJSON)
	if err != nil {
		return err
	}

	return nil
}
