package protocol

type FollowerProtocol struct {
}

func (f FollowerProtocol) Parse(input string) (Request, error) {
	return Request{}, nil
}
