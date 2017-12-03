package main

import (
	"fmt"
	"time"
	"xparser"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/transport"
)

type Processor struct {
	service   string
	pf        ProtocolFactory
	tf        TransportFactory
	tw        TransportWrapper
	chSuccess chan int
	chError   chan int32

	thrift_file string
	api_file    string
	case_name   string
	message     []byte
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
		if err := p.call(proto); err != nil {
			fmt.Println(gid, err)
			return
		}
		duration := time.Now().UnixNano() - snano
		p.chSuccess <- int(duration / 1000)
	}
}

func (p *Processor) test() (err error) {
	trans, err := p.GetTransport()
	if err != nil {
		return
	}
	defer trans.Close()
	proto := p.GetProtocol(trans)
	return p.call(proto)
}

func (p *Processor) call(proto Protocol) (err error) {
	var trans = proto.GetTransport()
	if _, err = trans.Write(p.message); err != nil {
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

func (p *Processor) flushToTrans(trans Transport, skip_result bool) (err error) {
	trans = p.tw.GetTransport(trans)
	proto := p.pf.GetProtocol(trans)

	if p.service != "" {
		proto = NewMultiplexedProtocol(proto, p.service)
	}

	// thrift parser
	parser_instance, err := xparser.InitParser(p.thrift_file)
	if err != nil {
		return
	}

	// api case
	var api_case *xparser.APICase
	if p.case_name == "" {
		if api_case, err = xparser.GetPingCase(); err != nil {
			return
		}
	} else if api_case, err = xparser.GetCase(p.api_file, p.case_name); err != nil {
		return
	}

	if err = parser_instance.BuildRequest(proto, api_case); err != nil {
		return
	}
	err = proto.GetTransport().Flush()

	return skip_response(proto)
}

func (p *Processor) initMessage() (err error) {
	membuffer := NewTMemoryBuffer()
	membuffer.Open()
	defer membuffer.Close()

	proto := p.GetProtocol(membuffer)

	var fn = call("ping")
	if err = fn(proto); err != nil {
		return
	}
	p.message = membuffer.GetBytes()
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
