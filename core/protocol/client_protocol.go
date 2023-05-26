package protocol

import (
	"fmt"
	"strings"
)

var (
	CMDGet  Command = Command{Op: "GET", MinRequiredParams: 1}
	CMDSet  Command = Command{Op: "SET", MinRequiredParams: 2}
	CMDDel  Command = Command{Op: "DEL", MinRequiredParams: 1}
	CMDPing Command = Command{Op: "PING", MinRequiredParams: 0}
)

type ClientProtocol struct {
}

func (f ClientProtocol) extractCommand(cmd string) (Command, error) {
	switch cmd {
	case CMDGet.Op:
		return CMDGet, nil
	case CMDSet.Op:
		return CMDSet, nil
	case CMDDel.Op:
		return CMDDel, nil
	case CMDPing.Op:
		return CMDPing, nil
	default:
		return Command{}, ErrInvalidCommand
	}
}

func (p ClientProtocol) Parse(input string) (Request, error) {
	input_parts := strings.Split(input, " ")
	fmt.Println(len(input_parts))
	if len(input_parts) == 0 {
		return Request{}, ErrInvalidRequest
	}

	command, err := p.extractCommand(input_parts[0])
	if err != nil {
		return Request{}, err
	}
	params := input_parts[1:]

	if !command.isParamsValid(params) {
		return Request{}, ErrInvalidRequest
	}

	return Request{Command: command, Params: params}, nil
}