package handlers

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/fsm"
	"log"
)

type Allocation struct {
	fsm *fsm.StateMachine

	status dispatcher.EventHandlerStatus
}

func NewAllocation() *Allocation {
	return &Allocation{status: dispatcher.Ready}
}

func (a *Allocation) Start(ctx context.Context) {
	a.fsm = fsm.New().WithContext(ctx)

	a.fsm.NewState().FromAny().To("Start").OnEnter(func(s *fsm.State, data interface{}) {
		a.status = dispatcher.Running
		log.Printf("\033[35m%s\033[0m\n", "AllocationHandler Start")

		a.fsm.Transition("Scan", nil)
	})

	// 开始到扫描
	a.fsm.NewState().FromAny().To("Scan").OnEnter(func(s *fsm.State, data interface{}) {
		master := a.fsm.CurrentState.Context().Value("dispather").(*dispatcher.EventDispathcer)

		eventQueue := master.EventQueue()
		if len(eventQueue) > 0 {
			event := eventQueue[0]
			master.Pop()
			a.fsm.Transition("Allocation", event)
			return
		}
		a.fsm.Transition("Scan", nil)
	})

	a.fsm.NewState().FromAny().To("Allocation").OnEnter(func(s *fsm.State, data interface{}) {
		master := a.fsm.CurrentState.Context().Value("dispather").(*dispatcher.EventDispathcer)

		event, ok := data.(dispatcher.Event)
		if !ok {
			log.Println("\033[1;32mevent is not dispatcher.Event\033[0m")
		}

		handlers := master.HandlerPool()[event.Type()]

		// 处理事件
		if len(handlers) > 0 {
			for _, handler := range handlers {
				if handler.Status() == dispatcher.Ready {
					log.Printf("\033[1;32m%s: %s handle %s \033[0m\n", "AllocationHandler Allocation", handler.Type(), event.Type())
					a.fsm.Transition("Scan", event)
					return
				}
			}
		}

		// 退还事件
		log.Printf("\033[1;32m%s: %s cant'be handled \033[0m\n", "AllocationHandler Allocation", event.Type())
		master.AddEvent(event)
		a.fsm.Transition("Scan", nil)
	})

	a.fsm.Start()
	a.fsm.Transition("Start", nil)
}

func (a *Allocation) Status() dispatcher.EventHandlerStatus {
	return a.status
}

func (a *Allocation) Type() dispatcher.EventType {
	return "Allocation"
}

func (a *Allocation) Handle(e dispatcher.Event) {
}
