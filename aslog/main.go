package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var cnt = 0

var logger *logrus.Logger

type write struct {
}

func (w write) Write(b []byte) (int, error) {
	if cnt == 0 {
		cnt++
		time.Sleep(time.Second * 4)
	}
	return os.Stdout.Write(b)
}

func testw() io.Writer {
	return &write{}
}

func printmsg(msg string, wg *sync.WaitGroup) {
	fmt.Println("appear", msg)
	logger.WithFields(logrus.Fields{
		"\nname":  "T",
		"\nstate": "",
	}).Warning("yo ", msg)
	// log.Println("yo", msg)
	// fmt.Fprintln(testw(), "yo", msg)

	fmtr := logrus.TextFormatter{}
	b, err := fmtr.Format(logger.WithFields(logrus.Fields{
		"name":  "T",
		"state": "",
	}))

	if err != nil {
		fmt.Println("fmtr err:", err)
	}

	fmt.Println("formatted b:", string(b))

	wg.Done()
}

func main() {

	logrus.RegisterExitHandler(func() {
		fmt.Println("exiting through handler")
	})

	logger = logrus.New()
	logger.SetNoLock()

	logger.SetOutput(testw())
	// logrus.SetOutput(testw())

	// log.SetOutput(testw())
	// log.set
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go printmsg("hagu", &wg)
	time.Sleep(time.Second * 2)
	go printmsg("padu", &wg)

	wg.Wait()
	logger.Fatal("testing handler")

}
