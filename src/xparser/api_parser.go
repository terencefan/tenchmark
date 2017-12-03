package xparser

import (
	"encoding/json"
	"io/ioutil"
)

type APIParser struct {
	cases map[string]*APICase
}

type APICase struct {
	Service  string                 `json:"service"`
	Function string                 `json:"function"`
	Args     map[string]interface{} `json:"args"`
}

var PingCase = &APICase{
	Function: "ping",
}

func NewApiParser(file string) (parser *APIParser, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	var cases map[string]*APICase
	if err = json.Unmarshal(data, &cases); err != nil {
		return
	}

	parser = &APIParser{
		cases: cases,
	}
	return
}

func (a *APIParser) GetCase(name string) *APICase {
	return a.cases[name]
}

func (a *APIParser) GetCases() map[string]*APICase {
	return a.cases
}
