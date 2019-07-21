package main

import (
	"fmt"
	"os"

	"github.com/iwanbk/gobeanstalk"
)

type PapaBeanstalk struct {
	ServerAddress    string
	serverConnection *gobeanstalk.Conn
}

func (papa *PapaBeanstalk) Connect() {
	beanstalkConnection, err := gobeanstalk.Dial(papa.ServerAddress)
	if err != nil {
		// do retries or whatever you need
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("connected!")
	papa.serverConnection = beanstalkConnection
}

func (papa *PapaBeanstalk) Close() {
	if papa.serverConnection != nil {
		papa.serverConnection.Quit()
	}
}
