package syslog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Length      int
	Priority    int
	Version     int
	Hostname    string
	Application string
	Created     time.Time
	Body        string
}

func ParseMessage(b []byte) (Message, error) {
	msg := strings.SplitN(string(b), " ", 7)
	if len(msg) != 7 {
		return Message{}, fmt.Errorf("Expected message to have 7 parts, not %d", len(msg))
	}
	timestamp, app, process, body := msg[1], msg[3], msg[4], msg[6]
	created, err := time.Parse("2006-01-02T15:04:05+00:00", timestamp)
	if err != nil {
		created = time.Time{}
	}
	return Message{
		Length:      len(b),
		Hostname:    app,
		Application: process,
		Body:        body,
		Created:     created,
	}, nil
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
