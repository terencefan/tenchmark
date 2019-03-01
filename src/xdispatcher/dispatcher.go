package xdispatcher

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

type Dispatcher interface {
	Init()
	GetCase() ([]byte, int)
}

type DispatchBase struct {
	api_case_buffers [][]byte
	dispatcher_type  int
}

const (
	SPECIFIC    = iota
	ROUND_ROBIN
	RANDOM
	WEIGHT
)

type SpecificDispatch struct {
	DispatchBase
	Dispatcher
}

func (s *SpecificDispatch) Init() {
}

func (s *SpecificDispatch) GetCase() ([]byte, int) {
	return s.api_case_buffers[0], 0
}

type RoundRobinDispatch struct {
	DispatchBase
	Dispatcher

	index int
	mutex sync.Mutex
}

func (r *RoundRobinDispatch) Init() {
	r.index = 0
}

func (r *RoundRobinDispatch) GetCase() ([]byte, int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	index := r.index
	buffer := r.api_case_buffers[index]
	r.index = (index + 1) % len(r.api_case_buffers)
	return buffer, index
}

type RandomDispatch struct {
	DispatchBase
	Dispatcher
}

func (r *RandomDispatch) Init() {
	rand.Seed(time.Now().Unix())
}

func (r *RandomDispatch) GetCase() ([]byte, int) {
	index := rand.Intn(len(r.api_case_buffers))
	return r.api_case_buffers[index], index
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
