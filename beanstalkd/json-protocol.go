package main

import "encoding/json"

type JsonCommentProtocol struct {
}

func (protocol *JsonCommentProtocol) Decode(encodedComment []byte) (*Comment, error) {
	unCodedComment := Comment{}
	err := json.Unmarshal(encodedComment, &unCodedComment)
	if err != nil {
		return nil, err
	}
	return &unCodedComment, nil
}

func (protocol *JsonCommentProtocol) Encode(comment *Comment) ([]byte, error) {
	encodedComment, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}
	return encodedComment, nil
}

func MakeJsonCommentProtocol() *JsonCommentProtocol {
	return &JsonCommentProtocol{}
}
