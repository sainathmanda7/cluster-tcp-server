package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	IP       string
	ServerNo int
	FD       int
	Payload  string
}

// turn struct into wire format
func (m Message) Encode() string {
	return fmt.Sprintf("%s|%d|%d|%s",
		m.IP,
		m.ServerNo,
		m.FD,
		m.Payload,
	)
}

// parse incoming string into Message struct
func Decoder(s string) (*Message, error) {
	parts := strings.Split(s, "|")

	if len(parts) != 4 {
		return nil, fmt.Errorf("bad message format, expected 4 fields got %d", len(parts))
	}

	// parse server number
	serverNo, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("serverNo not valid: %v", err)
	}

	// parse fd
	fd, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("fd not valid: %v", err)
	}

	// build message
	return &Message{
		IP:       parts[0],
		ServerNo: serverNo,
		FD:       fd,
		Payload:  parts[3], // this was wrong in your code, should be 3
	}, nil
}