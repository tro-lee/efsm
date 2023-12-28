package handlers

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/fsm"
	"log"
	"math/rand"
	"reflect"
	"time"
)

type Parser struct {
	fsm *fsm.StateMachine

	status dispatcher.EventHandlerStatus
}

func NewParser() *Parser {
	return &Parser{status: dispatcher.Init}
}

func (p *Parser) Start(ctx context.Context) {
	p.fsm = fsm.New().WithContext(ctx)

	// 开机等待事件
	p.fsm.NewState().FromAny().To("Start").OnEnter(func(s *fsm.State, data interface{}) {
		log.Printf("\033[36m%s: %s ready\033[0m\n", "Parser-Handler", reflect.TypeOf(p))
		p.status = dispatcher.Ready
	})

	p.fsm.NewState().FromAny().To("Error").OnEnter(func(s *fsm.State, data interface{}) {
		log.Println("\033[31m Parser-Handler Error\033[0m")
		p.status = dispatcher.Error

		master := p.fsm.CurrentState.Context().Value("dispather").(*dispatcher.EventDispathcer)
		master.AddEvent(data.(dispatcher.Event))

		p.fsm.Transition("Start", nil)
	})

	p.fsm.NewState().From("Start").To("Running").OnEnter(func(s *fsm.State, data interface{}) {
		data, ok := data.(dispatcher.Event)
		if !ok || data.(dispatcher.Event).Type() != "Parser" {
			p.fsm.Transition("Error", data)
			return
		}

		log.Printf("\033[36m%s: %s start handle %s \033[0m\n", "Parser-Handler", reflect.TypeOf(p), data.(dispatcher.Event).Type())
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
		log.Printf("\033[36m%s: %s end handle %s \033[0m\n", "Parser-Handler", reflect.TypeOf(p), data.(dispatcher.Event).Type())
		p.fsm.Transition("Start", nil)
	})

	p.fsm.Start()
	p.fsm.Transition("Start", nil)
}

func (p *Parser) Handle(e dispatcher.Event) {
	p.fsm.Transition("Running", e)
}

func (p *Parser) Status() dispatcher.EventHandlerStatus {
	return p.status
}

func (p *Parser) Type() dispatcher.EventType {
	return "Parser"
}

// 解析事件

type ParserEvent struct {
}

func (p *ParserEvent) Type() dispatcher.EventType {
	return "Parser"
}

func (p *ParserEvent) Data() []byte {
	return nil
}
