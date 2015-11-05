package syslog

import (
	"fmt"
	"strconv"
	"time"
)

type Message struct {
	Length      int
	Priority    int
	Version     int
	Hostname    string
	Application string
	Stamp       time.Time
	Body        string
}

func ParseMessage(b []byte) (Message, error) {
	return Message{Length: len(b)}, nil
}

func ParseFrame(b []byte, expected int) ([]Message, error) {
	msgs := []Message{}
	ci := 0 // Current index
	ms := 0

	for ci < len(b) {
		if b[ci] == 32 {
			l, err := strconv.ParseInt(string(b[ms:ci]), 10, 32)
			if err != nil {
				return msgs, err
			}

			start := ci + 1
			end := start + int(l)

			if end > len(b) {
				return msgs, fmt.Errorf("Yikes")
			}

			msg, err := ParseMessage(b[start:end])
			if err != nil {
				return msgs, err
			}
			msgs = append(msgs, msg)

			ci += int(l) + 1
			ms = ci
		} else {
			ci++
		}
	}

	if len(msgs) != expected {
		return msgs, fmt.Errorf("expected %d syslog messages, parsed %d", expected, len(msgs))
	}
	return msgs, nil
}
