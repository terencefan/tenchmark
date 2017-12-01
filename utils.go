package main

import (
	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

func sort(values []int, l, r int) {
	if l >= r {
		return
	}

	pivot := values[l]
	i := l + 1

	for j := l + 1; j <= r; j++ {
		if pivot > values[j] {
			values[i], values[j] = values[j], values[i]
			i++
		}
	}

	values[l], values[i-1] = values[i-1], pivot

	sort(values, l, i-2)
	sort(values, i, r)
}

func read_exception(proto Protocol) (err error) {
	var ae *TApplicationException
	if ae, err = ReadTApplicationException(proto); err != nil {
		return err
	}
	if err = proto.ReadMessageEnd(); err != nil {
		return err
	}
	return ae
}
