package xparser

import (
	"testing"

	"github.com/stdrickforce/thriftgo/protocol"
	"github.com/stdrickforce/thriftgo/thrift"
	"github.com/stdrickforce/thriftgo/transport"
)

func TestGetCallArgs(t *testing.T) {
	api_case, err := GetCase("../../example/api-RevenueOrder.json", "case1")
	if err != nil {
		panic(err)
	}

	//for _, arg := range args {
	//if arg.Type.Name == "list" {
	//fmt.Println(arg.Type.ValueType.Name)
	//fmt.Println(arg.Type.ValueType)
	//}
	////fmt.Println(arg.Type)
	//}
	p, err := InitParser("../../example/RevenueOrder.thrift")
	if err != nil {
		panic(err)
	}
	//trans := transport.NewTMemoryBuffer()
	var trans transport.Transport
	trans = transport.NewTSocket("127.0.0.1:8888")
	trans.Open()
	trans = transport.NewTBufferedTransport(trans)
	proto := protocol.NewTBinaryProtocol(trans, true, true)

	err = p.BuildRequest(proto, api_case)

	if err != nil {
		panic(err)
	}

	proto.Flush()
	proto.ReadMessageBegin()
	proto.Skip(thrift.T_STRUCT)
	proto.ReadMessageEnd()
	trans.Close()

	//fmt.Println(trans.GetBytes())

	//fmt.Println(args)
}
