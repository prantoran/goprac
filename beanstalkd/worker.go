package main

import (
	"fmt"
	"time"

	"github.com/iwanbk/gobeanstalk"
	"github.com/sirupsen/logrus"
)

type CommentWorker struct {
	PapaBeanstalk
	protocol            CommentProtocol
	processor           CommentProcessor
	reserveWaitDuration int // seconds
}

func (worker *CommentWorker) ProcessJob(tubeName string) {
	fmt.Println("reserving job")

	if tubecnt, err := worker.serverConnection.Watch(tubeName); err != nil {
		logrus.Error(err)
		return
	} else {
		logrus.Info("no. of tubes watching:", tubecnt)
	}

	job, err := worker.serverConnection.Reserve(time.Duration(worker.reserveWaitDuration) * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("got job Id: ", job.ID)
	comment, err := worker.protocol.Decode(job.Body)
	if err != nil {
		worker.handleError(job, err)
		return
	}
	processError := worker.processor.DoProcess(comment)
	if processError != nil {
		worker.handleError(job, err)
		return
	}
	fmt.Println("processed job id: ", job.ID)
	worker.serverConnection.Delete(job.ID)
}

func (worker *CommentWorker) handleError(job *gobeanstalk.Job, err error) {
	fmt.Println(err)
	if job == nil {
		return
	}
	priority := uint32(5)
	delay := 0 * time.Second
	worker.serverConnection.Release(job.ID, priority, delay) // hey I can't handle this
}

func MakeNewWorker(serverAddress string, protocol CommentProtocol, processor CommentProcessor) *CommentWorker {
	worker := CommentWorker{protocol: protocol, processor: processor}
	worker.ServerAddress = serverAddress
	worker.reserveWaitDuration = 100
	return &worker
}
