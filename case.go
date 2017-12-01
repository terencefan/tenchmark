package main

import (
	"errors"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type Case func(proto Protocol) error

func call(name string, args ...interface{}) Case {
	var writeMessageBody = func(proto Protocol) (err error) {
		if err = proto.WriteStructBegin("whatever"); err != nil {
			return
		}
		for i, arg := range args {
			index := int16(i + 1)
			switch v := arg.(type) {
			case int16:
				err = proto.WriteFieldBegin("i16", T_I16, index)
				err = proto.WriteI16(v)
			case int32:
				err = proto.WriteFieldBegin("i32", T_I32, index)
				err = proto.WriteI32(v)
			case int64:
				err = proto.WriteFieldBegin("i64", T_I64, index)
				err = proto.WriteI64(v)
			case string:
				err = proto.WriteFieldBegin("string", T_STRING, index)
				err = proto.WriteString(v)
			default:
				err = errors.New("unsupport type")
			}
			if err != nil {
				return
			}
		}
		if err = proto.WriteFieldStop(); err != nil {
			return
		}
		if err = proto.WriteStructEnd(); err != nil {
			return
		}
		if err = proto.WriteMessageEnd(); err != nil {
			return
		}
		if err = proto.Flush(); err != nil {
			return
		}
		return
	}

	return func(proto Protocol) (err error) {
		if err = proto.WriteMessageBegin(name, T_CALL, 0); err != nil {
			return
		}
		if err = writeMessageBody(proto); err != nil {
			return
		}
		_, mtype, _, err := proto.ReadMessageBegin()
		if err != nil {
			return
		} else if mtype == T_EXCEPTION {
			return read_exception(proto)
		}
		if err = proto.Skip(T_STRUCT); err != nil {
			return
		}
		if err = proto.ReadMessageEnd(); err != nil {
			return
		}
		return
	}
}

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
