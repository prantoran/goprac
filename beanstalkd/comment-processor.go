package main

import (
	"fmt"
	"os"
	"time"
)

type DiskCommentProcessor struct {
	dir string
}

func (processor *DiskCommentProcessor) DoProcess(comment *Comment) error {
	filePath := fmt.Sprintf("%s/%s_%s.txt", processor.dir, comment.Date.Format(time.RFC3339), comment.UserName)
	fmt.Println("DoProcess() filePath:", filePath)
	commentFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer commentFile.Close()
	commentFile.Write([]byte(comment.Text))
	return nil
}

func MakeNewCommentProcessor(dir string) *DiskCommentProcessor {
	return &DiskCommentProcessor{dir: dir}
}
