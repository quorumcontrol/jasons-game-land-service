package service

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
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

func (sa *ServiceActor) Receive(context actor.Context) {
	msg := context.Message()
	switch msg.(type) {
	case *actor.Started:
		fmt.Printf("Service actor started\n")
	case *actor.Stopped:
		fmt.Printf("Service actor stopped\n")
	}
}
