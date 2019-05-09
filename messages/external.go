//go:generate msgp

package messages

import (
	sdkmessages "github.com/quorumcontrol/tupelo-go-sdk/gossip3/messages"
)

func init() {
	sdkmessages.RegisterMessage(&BuildPortal{})
	sdkmessages.RegisterMessage(&BuildPortalResponse{})
}

// BuildPortal represents a request to build a portal.
type BuildPortal struct {
}

func (BuildPortal) TypeCode() int8 {
	return 1
}

// BuildPortalResponse represents a response to a request to build a portal.
type BuildPortalResponse struct {
}

func (BuildPortalResponse) TypeCode() int8 {
	return 2
}
