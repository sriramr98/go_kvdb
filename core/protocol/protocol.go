package protocol

import (
	"errors"
)

var ErrInvalidCommand = errors.New("invalid command")
var ErrInvalidRequest = errors.New("invalid request")

type Request struct {
	Command Command
	Params  []string
}

type Response struct {
	Success bool
	Value   []byte
}
type Protocol interface {
	Parse(input string) (Request, error)
}

type Command struct {
	Op                 string
	MinRequiredParams  int
	IsReplicable       bool
	CanFollowerProcess bool
}

func (c Command) isParamsValid(params []string) bool {
	return len(params) >= c.MinRequiredParams
}
