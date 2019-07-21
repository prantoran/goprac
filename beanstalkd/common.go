package main

import "time"

type Comment struct {
	UserName string
	Date     time.Time
	Text     string
}

type CommentProtocol interface {
	Decode(encodedComment []byte) (*Comment, error)
	Encode(comment *Comment) ([]byte, error)
}

type CommentProcessor interface {
	DoProcess(comment *Comment) error
}
