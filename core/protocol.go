package core

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid command")

type Request struct {
	Command Command
	Params  []string
}

type Response struct {
	Value   []byte
	Success bool
}

type Command struct {
	Cmd               string
	MinRequiredParams int
}

var (
	CMDGet  Command = Command{Cmd: "GET", MinRequiredParams: 1}
	CMDSet  Command = Command{Cmd: "SET", MinRequiredParams: 2}
	CMDDel  Command = Command{Cmd: "DEL", MinRequiredParams: 1}
	CMDPing Command = Command{Cmd: "PING", MinRequiredParams: 0}
)

func (c Command) isParamsValid(params []string) bool {
	return len(params) >= c.MinRequiredParams
}

func parseCommand(cmd string) (Command, error) {
	switch cmd {
	case CMDGet.Cmd:
		return CMDGet, nil
	case CMDSet.Cmd:
		return CMDSet, nil
	case CMDDel.Cmd:
		return CMDDel, nil
	case CMDPing.Cmd:
		return CMDPing, nil
	default:
		return Command{}, ErrInvalidCommand
	}
}

func ParseProtocol(input string) (Request, error) {
	input_parts := strings.Split(input, " ")
	fmt.Println(len(input_parts))
	if len(input_parts) == 0 {
		return Request{}, errors.New("empty request")
	}

	command, err := parseCommand(input_parts[0])
	if err != nil {
		return Request{}, err
	}
	params := input_parts[1:]

	// TODO: Add better validation - For example, if EX is set, then the value should be present and should be an integer for PUT command
	if !command.isParamsValid(params) {
		return Request{}, errors.New("invalid params")
	}

	return Request{Command: command, Params: params}, nil
}
