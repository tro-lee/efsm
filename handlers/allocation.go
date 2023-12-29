package handlers

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/fsm"
	"learn/efsm/out"
	"log"
	"reflect"
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
		out.Start("%s Start", reflect.TypeOf(a))

		a.fsm.Transition("Scan", nil)
	})

	// 开始到扫描
	a.fsm.NewState().FromAny().To("Scan").OnEnter(func(s *fsm.State, data interface{}) {
		master := a.fsm.CurrentState.Context().Value("dispather").(*dispatcher.EventDispathcer)

		event := <-master.EventQueue()
		a.fsm.Transition("Allocation", event)
	})

	a.fsm.NewState().FromAny().To("Allocation").OnEnter(func(s *fsm.State, data interface{}) {
		master := a.fsm.CurrentState.Context().Value("dispather").(*dispatcher.EventDispathcer)

		event := data.(dispatcher.Event)
		handlers := master.HandlerPool()[event.Type()]

		// 处理事件
		if len(handlers) > 0 {
			for _, handler := range handlers {
				if handler.Status() == dispatcher.Ready {
					handler.SetStatus(dispatcher.Running)
					log.Printf("\033[1;35m%s: %s handle %s \033[0m\n", "AllocationHandler Allocation", reflect.TypeOf(handler), event.Type())
					handler.Handle(event)
					a.fsm.Transition("Scan", event)
					return
				}
			}
		}

		// 退还事件
		master.AddEvent(event)
		a.fsm.Transition("Scan", nil)
	})

	a.fsm.Start()
	a.fsm.Transition("Start", nil)
}

func (a *Allocation) Status() dispatcher.EventHandlerStatus {
	return a.status
}

func (a *Allocation) SetStatus(status dispatcher.EventHandlerStatus) {
	a.status = status
}

func (a *Allocation) Type() dispatcher.EventType {
	return "Allocation"
}

func (a *Allocation) Handle(e dispatcher.Event) {
}
