package service

import (
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
	case *sdkmessages.Ping:
		ctx.Respond(&sdkmessages.Pong{Msg: msg.Msg})
	case *messages.BuildPortal:
		ctx.Respond(&messages.BuildPortalResponse{})
	}
}
