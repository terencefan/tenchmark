package xdispatcher

import (
	"fmt"
	"math/rand"
	"time"
)

type Dispatcher interface {
	Init()
	GetCase() []byte
}

type DispatchBase struct {
	api_case_buffers [][]byte
	dispatcher_type  int
}

const (
	SPECIFIC = iota
	ROUND_ROBIN
	RANDOM
	WEIGHT
)

type SpecificDispatch struct {
	DispatchBase
	Dispatcher
}

func (s *SpecificDispatch) GetCase() []byte {
	return s.api_case_buffers[0]
}

type RoundRobinDispatch struct {
	DispatchBase
	Dispatcher

	index int
}

func (r *RoundRobinDispatch) Init() {
	r.index = 0
}

func (r *RoundRobinDispatch) GetCase() []byte {
	buffer := r.api_case_buffers[r.index]
	r.index = (r.index + 1) % len(r.api_case_buffers)
	return buffer
}

type RandomDispatch struct {
	DispatchBase
	Dispatcher
}

func (r *RandomDispatch) Init() {
	rand.Seed(time.Now().Unix())
}

func (r *RandomDispatch) GetCase() []byte {
	index := rand.Intn(len(r.api_case_buffers))
	return r.api_case_buffers[index]
}

func NewDispatcher(api_case_buffers [][]byte, dispatcher_type int) (dispatcher Dispatcher, err error) {
	switch dispatcher_type {
	case SPECIFIC:
		dispatcher = &SpecificDispatch{DispatchBase: DispatchBase{
			api_case_buffers: api_case_buffers,
			dispatcher_type:  dispatcher_type,
		}}
	case ROUND_ROBIN:
		dispatcher = &RoundRobinDispatch{DispatchBase: DispatchBase{
			api_case_buffers: api_case_buffers,
			dispatcher_type:  dispatcher_type,
		}}
	case RANDOM:
		dispatcher = &RandomDispatch{DispatchBase: DispatchBase{
			api_case_buffers: api_case_buffers,
			dispatcher_type:  dispatcher_type,
		}}
	case WEIGHT:
		fmt.Errorf("unsupport type \"WEIGHT\" of dispatcher")
	default:
		fmt.Errorf("unknown type of dispatcher")
	}

	dispatcher.Init()
	return
}
