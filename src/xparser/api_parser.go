package xparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type APICase struct {
	Service  string
	Function string
	Args     map[string]interface{}
}

func GetPingCase() (api_case *APICase, err error) {
	api_case = &APICase{
		Function: "ping",
	}
	return
}

func GetCase(file string, case_name string) (api_case *APICase, err error) {
	var (
		data     []byte
		obj      interface{}
		api_map  map[string]interface{}
		case_map map[string]interface{}
		ok       bool
	)

	api_case = new(APICase)

	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		return
	}

	if api_map, ok = obj.(map[string]interface{}); !ok {
		err = fmt.Errorf("json data cannot convert to map: %v", data)
		return
	}
	if test_case, ok := api_map[case_name]; ok {
		if case_map, ok = test_case.(map[string]interface{}); !ok {
			err = fmt.Errorf("json data cannot convert to map: %v", data)
			return
		}

		if api_case.Service, ok = case_map["service"].(string); !ok {
			err = fmt.Errorf("service non-exist or invalid service in case")
			return
		}

		if api_case.Function, ok = case_map["function"].(string); !ok {
			err = fmt.Errorf("function non-exist or invalid function in case")
			return
		}

		if api_case.Args, ok = case_map["args"].(map[string]interface{}); !ok {
			err = fmt.Errorf("args non-exist or invalid args in case")
			return
		}
	} else {
		err = fmt.Errorf("api file have no case named: %s")
		return
	}
	return
}
