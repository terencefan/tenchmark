package main

import (
	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

func create_order_body(proto Protocol) error {
	proto.WriteStructBegin("whatever")

	// appid
	proto.WriteFieldBegin("appId", T_I16, 1)
	proto.WriteI16(2)
	proto.WriteFieldEnd()

	// torder
	proto.WriteFieldBegin("tOrder", T_STRUCT, 2)
	proto.WriteStructBegin("TOrder")

	// torder.siteid
	proto.WriteFieldBegin("siteId", T_I32, 3)
	proto.WriteI32(2)
	proto.WriteFieldEnd()

	// torder.uid
	proto.WriteFieldBegin("uid", T_I32, 5)
	proto.WriteI32(2)
	proto.WriteFieldEnd()

	// torder.type
	proto.WriteFieldBegin("type", T_I32, 7)
	proto.WriteI32(2013)
	proto.WriteFieldEnd()

	// torder.price
	proto.WriteFieldBegin("price", T_I32, 9)
	proto.WriteI32(10000)
	proto.WriteFieldEnd()

	// torder stop
	proto.WriteFieldStop()

	proto.WriteFieldEnd()
	proto.WriteFieldStop()

	proto.WriteStructEnd()
	proto.WriteMessageEnd()
	proto.Flush()
	return nil
}

func create_order(proto Protocol) (err error) {
	if err = proto.WriteMessageBegin("RevenueOrder:createOrder", T_CALL, 0); err != nil {
		return
	}
	if err = create_order_body(proto); err != nil {
		return
	}
	if _, _, _, err = proto.ReadMessageBegin(); err != nil {
		return
	}
	if err = proto.Skip(T_STRUCT); err != nil {
		return
	}
	if err = proto.ReadMessageEnd(); err != nil {
		return
	}
	return
}
