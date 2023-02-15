package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	//kafka.Message{
	//	Topic:         "",
	//	Partition:     0,
	//	Offset:        0,
	//	HighWaterMark: 0,
	//	Key:           nil,
	//	Value:         nil,
	//	Headers: []kafka.Header{
	//		{
	//			"name",
	//			[]byte("hahaha"),
	//		},
	//		{},
	//	},
	//	Time: time.Time{},
	//}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
