package xparser

import (
	"errors"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type Case func(proto Protocol) error

func Call(name string, args ...interface{}) Case {
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
		if err = proto.GetTransport().Flush(); err != nil {
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
		return
	}
}
