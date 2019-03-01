package main

import (
	"fmt"
	"time"
	"xdispatcher"
	"xparser"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/transport"
)

type Processor struct {
	service   string
	pf        ProtocolFactory
	tf        TransportFactory
	tw        TransportWrapper
	chSuccess chan *CallInfo
	chError   chan int32

	dispatcher xdispatcher.Dispatcher
}

func (p *Processor) process(gid int, pipe <-chan int) {
	defer wg.Done()

	trans := p.tf.GetTransport()
	if err := trans.Open(); err != nil {
		panic(err)
	}
	defer trans.Close()

	proto := p.pf.GetProtocol(trans)

	for _ = range pipe {
		snano := time.Now().UnixNano()
		buffers, index := p.dispatcher.GetCase()
		if err := p.call(proto, buffers); err != nil {
			fmt.Println(gid, err)
			return
		}
		duration := time.Now().UnixNano() - snano
		p.chSuccess <- &CallInfo{index: index, duration: int(duration / 1000)}
	}
}

func (p *Processor) test() (err error) {
	trans, err := p.GetTransport()
	trans.SetTimeout(5 * time.Second)
	if err != nil {
		return
	}
	defer trans.Close()
	proto := p.GetProtocol(trans)
	buffers, _ := p.dispatcher.GetCase()
	return p.call(proto, buffers)
}

func (p *Processor) call(proto Protocol, buffers []byte) (err error) {
	var trans = proto.GetTransport()
	if _, err = trans.Write(buffers); err != nil {
		return
	}
	if err = trans.Flush(); err != nil {
		return
	}
	if err = skip_response(proto); err != nil {
		return
	}
	return
}

func (p *Processor) GetTransport() (trans Transport, err error) {
	trans = p.tf.GetTransport()
	trans = p.tw.GetTransport(trans)
	if err = trans.Open(); err != nil {
		return nil, err
	}
	return trans, nil
}

func (p *Processor) GetProtocol(trans Transport) Protocol {
	var proto = p.pf.GetProtocol(trans)
	if p.service != "" {
		proto = NewMultiplexedProtocol(proto, p.service)
	}
	return proto
}

func InitRunDispatcher(filename string, processor *Processor) (err error) {
	var (
		loader xdispatcher.DataLoader
	)
	if filename != "" {
		if loader, err = xdispatcher.NewFileDataLoader(filename); err != nil {
			return
		}
		loader.Load()
		if processor.dispatcher, err = xdispatcher.NewDispatcher(loader.GetAllApis(), xdispatcher.SPECIFIC); err != nil {
			return
		}
	} else {
		trans := NewTMemoryBuffer()
		proto := processor.GetProtocol(trans)
		xparser.BuildPing(proto, xparser.PingCase)
		apis := [][]byte{trans.GetBytes()}
		if processor.dispatcher, err = xdispatcher.NewDispatcher(apis, xdispatcher.SPECIFIC); err != nil {
			return
		}
	}
	return
}
