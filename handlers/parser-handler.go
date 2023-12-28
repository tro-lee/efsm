package handlers

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/fsm"
	"log"
)

type Parser struct {
	fsm *fsm.StateMachine

	status dispatcher.EventHandlerStatus
}

func NewParser() *Parser {
	return &Parser{status: dispatcher.Init}
}

func (p *Parser) Start(ctx context.Context) {
	p.fsm.WithContext(ctx)

	// 开机等待事件
	p.fsm.NewState().FromAny().To("Start").OnEnter(func(s *fsm.State, data interface{}) {
		if s.Destination == "Error" {
			return
		}

		p.status = dispatcher.Ready
	})
	p.fsm.NewState().FromAny().To("Error").OnEnter(func(s *fsm.State, data interface{}) {

	})
	p.fsm.NewState().From("Start").To("Running").OnEnter(func(s *fsm.State, data interface{}) {
		data, ok := data.(dispatcher.Event)
		if !ok || data.(dispatcher.Event).Type() != "Parser" {
			log.Println("\033[1;31m Parser-Handler Error\033[0m")
		}
	})

}
