package main

import (
	"fmt"
	"time"
)

type Producer struct {
	PapaBeanstalk
	protocol CommentProtocol
}

func (producer *Producer) PutComment(comment *Comment, tubeName string) error {

	if err := producer.serverConnection.Use(tubeName); err != nil {
		return err
	}

	body, err := producer.protocol.Encode(comment)
	if err != nil {
		return err
	}
	priority := uint32(10)
	delay := 0 * time.Second
	time_to_run := 20 * time.Second

	jobId, err := producer.serverConnection.Put(body, priority, delay, time_to_run)
	if err != nil {
		return err
	}
	fmt.Println("inserted Job id: ", jobId)
	return nil
}

func MakeNewProducer(serverAdress string, protocol CommentProtocol) *Producer {
	producer := Producer{protocol: protocol}
	producer.ServerAddress = serverAdress
	return &producer
}
