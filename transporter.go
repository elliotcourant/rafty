package rafty

import (
	"io"
)

type RPCResponse struct {
	Response interface{}
	Error    error
}

type RPC struct {
	Command         interface{}
	Reader          io.Reader // Only when installing a snapshot.
	ResponseChannel chan<- RPCResponse
}

func (r *RPC) Respond(response interface{}, err error) {
	r.ResponseChannel <- RPCResponse{response, err}
}
