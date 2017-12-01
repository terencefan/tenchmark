package xparser

import (
	"fmt"
	"strconv"

	"github.com/stdrickforce/go-thrift/parser"
	"github.com/stdrickforce/thriftgo/protocol"
	"github.com/stdrickforce/thriftgo/thrift"
)

type ThriftParser struct {
	Thrift_name string
	Th          *thrift.Thrift
}

type FieldMeta struct {
	IsField bool
	ID      int16
	Name    string
	Type    *parser.Type
}

var (
	parser_instance *ThriftParser
)

func InitParser(thrift_file string) (p *ThriftParser, err error) {
	if parser_instance == nil {
		parser_instance = new(ThriftParser)
	}
	parser_instance.Thrift_name = thrift_file
	parser_instance.Th, _, err = thrift.Parse(thrift_file)
	p = parser_instance
	return
}

func GetParser() (*ThriftParser, error) {
	if parser_instance == nil {
		return nil, fmt.Errorf("parser root is nil")
	}
	return parser_instance, nil

}

func (p *ThriftParser) GetCallArgs(api_case *APICase) (args []*parser.Field, err error) {
	var (
		service *parser.Service
		ok      bool
	)
	if service, ok = p.Th.Services[api_case.Service]; !ok {
		if len(p.Th.Services) > 0 {
			for key := range p.Th.Services {
				service = p.Th.Services[key]
				break
			}
		} else {
			err = fmt.Errorf("thrift file have no service named: %s", api_case.Service)
			return
		}
	}

	if function, ok := service.Methods[api_case.Function]; ok {
		args = function.Arguments
	} else {
		err = fmt.Errorf("service %s have no function named: %s", api_case.Service, api_case.Function)
	}
	return
}

func (p *ThriftParser) GetStruct(struct_name string) (st *parser.Struct, err error) {
	var ok bool
	if p.Th == nil {
		err = fmt.Errorf("no thrift parser root, struct: %s", struct_name)
		return
	}
	if st, ok = p.Th.Structs[struct_name]; !ok {
		err = fmt.Errorf("no such struct in thrift file, struct: %s", struct_name)
		return
	}
	return
}

func (p *ThriftParser) BuildRequest(proto protocol.Protocol, api_case *APICase) (err error) {
	var args []*parser.Field

	if args, err = p.GetCallArgs(api_case); err != nil {
		return
	}

	if err = proto.WriteMessageBegin(api_case.Function, thrift.T_CALL, 0); err != nil {
		return
	}
	if err = proto.WriteStructBegin(""); err != nil {
		return
	}
	jsonData := api_case.Args

	for _, arg := range args {
		var (
			data interface{}
			ok   bool
		)
		if data, ok = jsonData[arg.Name]; ok {
			err = p.writeData(proto, data, &FieldMeta{
				IsField: true,
				ID:      int16(arg.ID),
				Name:    arg.Name,
				Type:    arg.Type,
			})
			if err != nil {
				return
			}
		} else if !arg.Optional {
			panic(fmt.Sprintf("arg %s is required!", arg.Name))
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

	return
}

func (p *ThriftParser) writeData(proto protocol.Protocol, data interface{}, field_meta *FieldMeta) (err error) {
	elem_type := TypeToByte(field_meta.Type.Name)
	if field_meta.IsField {
		proto.WriteFieldBegin(field_meta.Name, elem_type, field_meta.ID)
	}
	switch elem_type {
	case thrift.T_I16:
		var int_data int
		if float_data, ok := data.(float64); ok {
			int_data = int(float_data)
		} else {
			int_data, err = strconv.Atoi(data.(string))
		}
		if err != nil {
			return
		}
		err = proto.WriteI16(int16(int_data))
	case thrift.T_I32:
		var int_data int
		if float_data, ok := data.(float64); ok {
			int_data = int(float_data)
		} else {
			int_data, err = strconv.Atoi(data.(string))
		}
		if err != nil {
			return
		}
		err = proto.WriteI32(int32(int_data))
	case thrift.T_I64:
		var int_data int64
		if float_data, ok := data.(float64); ok {
			int_data = int64(float_data)
		} else {
			int_data, err = strconv.ParseInt(data.(string), 10, 64)
		}
		if err != nil {
			return
		}
		err = proto.WriteI64(int64(int_data))
	case thrift.T_BYTE:
		if byte_data, ok := data.(byte); ok {
			err = proto.WriteByte(byte_data)
		} else if str_data, ok := data.(string); ok {
			if len(str_data) > 0 {
				err = proto.WriteByte(byte(str_data[0]))
			} else {
				// NOTE 可能有问题
				err = proto.WriteByte(byte(0))
			}
		}
	case thrift.T_BOOL:
		var (
			bool_data bool
			ok        bool
		)
		if bool_data, ok = data.(bool); ok {
			err = proto.WriteBool(bool_data)
		} else {
			if bool_data, err = strconv.ParseBool(data.(string)); err != nil {
				return
			}
			err = proto.WriteBool(bool_data)
		}
	case thrift.T_DOUBLE:
		var (
			float_data float64
			ok         bool
		)
		if float_data, ok = data.(float64); ok {
			err = proto.WriteDouble(float_data)
		} else {
			if float_data, err = strconv.ParseFloat(data.(string), 64); err != nil {
				return
			}
			err = proto.WriteDouble(float_data)
		}
	case thrift.T_STRING:
		err = proto.WriteString(data.(string))
	case thrift.T_LIST:
		if val, ok := data.([]interface{}); ok {
			value_type := field_meta.Type.ValueType
			err = proto.WriteListBegin(TypeToByte(value_type.Name), len(val))
			for _, elem := range val {
				if err = p.writeData(proto, elem, &FieldMeta{IsField: false, Type: value_type}); err != nil {
					return
				}
			}
			err = proto.WriteListEnd()
		}
	case thrift.T_SET:
		if val, ok := data.([]interface{}); ok {
			value_type := field_meta.Type.ValueType
			err = proto.WriteSetBegin(TypeToByte(value_type.Name), len(val))
			for _, elem := range val {
				if err = p.writeData(proto, elem, &FieldMeta{IsField: false, Type: value_type}); err != nil {
					return
				}
			}
			err = proto.WriteSetEnd()
		}
	case thrift.T_MAP:
		if val, ok := data.(map[string]interface{}); ok {
			value_type := field_meta.Type.ValueType
			key_type := field_meta.Type.KeyType
			err = proto.WriteMapBegin(TypeToByte(key_type.Name), TypeToByte(value_type.Name), len(val))

			for key, elem := range val {
				err = p.writeData(proto, key, &FieldMeta{IsField: false, Type: key_type})
				err = p.writeData(proto, elem, &FieldMeta{IsField: false, Type: value_type})
				if err != nil {
					return
				}
			}
			err = proto.WriteMapEnd()
			err = proto.WriteFieldEnd()
		}
	case thrift.T_STRUCT:
		if val, ok := data.(map[string]interface{}); ok {
			var (
				p          *ThriftParser
				cur_struct *parser.Struct
			)
			if p, err = GetParser(); err != nil {
				return
			}
			if cur_struct, err = p.GetStruct(field_meta.Type.Name); err != nil {
				return
			}
			for _, field := range cur_struct.Fields {
				if field_data, ok := val[field.Name]; ok {
					err = p.writeData(proto, field_data, &FieldMeta{
						IsField: true,
						ID:      int16(field.ID),
						Name:    field.Name,
						Type:    field.Type,
					})
					if err != nil {
						return
					}
				} else if !field.Optional {
					err = fmt.Errorf("field %s is required!", field.Name)
					return
				}
			}
			proto.WriteFieldStop()
		}
	}

	//if field_meta.IsField {
	//proto.WriteFieldStop()
	//}
	return
}

func TypeToByte(name string) byte {
	switch name {
	case "i16":
		return thrift.T_I16
	case "i32":
		return thrift.T_I32
	case "i64":
		return thrift.T_I64
	case "byte":
		return thrift.T_BYTE
	case "bool":
		return thrift.T_BOOL
	case "double":
		return thrift.T_DOUBLE
	case "string":
		return thrift.T_STRING
	case "list":
		return thrift.T_LIST
	case "set":
		return thrift.T_SET
	case "map":
		return thrift.T_MAP
	default:
		return thrift.T_STRUCT
	}
}
