package main

import (
	"fmt"
	"os"
	"time"
)

func ProducerMain() {
	var comments = []Comment{
		{UserName: "some_user", Text: "i love your cat", Date: time.Now()},
		{UserName: "some_other_user", Text: "i prefer dogs", Date: time.Now()},
		{UserName: "another_user", Text: "please close this thread", Date: time.Now()},
		{UserName: "admin", Text: "thread closed - not relevant", Date: time.Now()},
	}
	protocol := MakeJsonCommentProtocol()
	producer := MakeNewProducer("localhost:11300", protocol)
	producer.Connect()
	defer producer.Close()

	for i, comment := range comments {
		if (i+1)%2 == 1 {
			producer.PutComment(&comment, "odd")
		} else {
			producer.PutComment(&comment, "even")
		}
	}
}

func WorkerMain(tubename string) {
	protocol := MakeJsonCommentProtocol()
	commentsDir := "./comments"
	os.Mkdir(commentsDir, 0777)
	processor := MakeNewCommentProcessor(commentsDir)
	worker := MakeNewWorker("localhost:11300", protocol, processor)
	worker.Connect()
	defer worker.Close()
	for {
		worker.ProcessJob(tubename)
	}

}

func printUsage() {
	fmt.Println("Usage: example-app worker/producer")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}
	if os.Args[1] == "worker" {
		if len(os.Args) < 3 {
			fmt.Println("Usage example-app worker tubename")
			os.Exit(1)
		}
		WorkerMain(os.Args[2])

	} else if os.Args[1] == "producer" {
		ProducerMain()
	} else {
		printUsage()
	}
}
