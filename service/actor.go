package service

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/quorumcontrol/jasons-game-land-service/messages"
	sdkmessages "github.com/quorumcontrol/tupelo-go-sdk/gossip3/messages"
)

type ServiceConfig struct {
}

type ServiceActor struct {
	Config ServiceConfig
}

func NewServiceActor(cfg ServiceConfig) *ServiceActor {
	return &ServiceActor{
		Config: cfg,
	}
}

func (sa *ServiceActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Printf("Service actor started\n")
	case *sdkmessages.Ping:
		fmt.Printf("Ping!\n")
		ctx.Respond(&sdkmessages.Pong{Msg: msg.Msg})
	case *messages.BuildPortal:
		fmt.Printf("Asked to build a portal!\n")
		ctx.Respond(&messages.BuildPortalResponse{})
	case *actor.Stopped:
		fmt.Printf("Service actor stopped\n")
	}
}
