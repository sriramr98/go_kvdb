package protocol

import "strings"

var (
	CMDSync Command = Command{Op: "SYNC", MinRequiredParams: 0}
)

type FollowerProtocol struct {
}

func (f FollowerProtocol) extractCommand(input string) (Command, error) {
	switch input {
	case CMDSync.Op:
		return CMDSync, nil
	default:
		return Command{}, ErrInvalidCommand
	}
}

func (f FollowerProtocol) Parse(input string) (Request, error) {
	input_parts := strings.Split(input, " ")
	if len(input_parts) == 0 {
		return Request{}, ErrInvalidRequest
	}

	command, err := f.extractCommand(input_parts[0])
	if err != nil {
		return Request{}, err
	}

	return Request{Command: command}, nil

}
