package main

import (
	"fmt"
	"context"
	"strconv"
	"time"
	"log"
	"os"
	kafka "github.com/segmentio/kafka-go"
)

const (
	topic 		     = "message-log"
	broker1Address = "localhost:29092"
	broker2Address = "localhost:39092"
)

func produce(ctx context.Context) {
	// initialize a counter
	i := 0

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address},
		Topic: topic,
		// wait until we get 10 messages before writing
		BatchSize: 10,
		// no matter what happens, write all pending messages
		// every 2 seconds
		BatchTimeout: 2 * time.Second,
		Logger: l,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
		})

		if err != nil {
			panic("could not write message " + err.Error())
		}

		fmt.Println("writes: ", i)
		i++

		time.Sleep(time.Second)
	}
}

func consume(ctx context.Context, groupID string, startOffset string) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages

	offset := kafka.FirstOffset
	if startOffset == "last" {
		offset = kafka.LastOffset
	}

	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address, broker2Address},
		Topic: topic,
		GroupID: groupID,
		MinBytes: 5,
		// the kafka library requires you to set the MaxBytes
		// in case the MinBytes are set
		MaxBytes: 1e6,
		// wait for at most 3 seconds before receiving new data
		MaxWait: 3 * time.Second,
		StartOffset: offset,
		Logger: l,
	})

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("couldn't read message " + err.Error())
		}

		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func main() {
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	go produce(ctx)
	// consume(ctx, "my-group", "first")
	consume(ctx, "first-offset-group", "first")
	// consume(ctx, "last-offset-group", "last")
}
