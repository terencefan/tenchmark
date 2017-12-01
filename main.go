package main

import (
	"fmt"
	"runtime"
	"sync"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/transport"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	wg sync.WaitGroup
)

var (
	requests          = kingpin.Flag("requests", "Number of requests to perform").Short('n').Default("1000").Int()
	concurrency       = kingpin.Flag("concurrency", "Number of multiple requests to make at a time").Short('c').Default("10").Int()
	protocol          = kingpin.Flag("protocol", "Specify protocol factory").Default("binary").String()
	transport         = kingpin.Flag("transport", "Specify transport factory").Default("socket").String()
	transport_wrapper = kingpin.Flag("transport-wrapper", "Specify transport wrapper").Default("buffered").String()
	service           = kingpin.Flag("service", "Specify service name").String()

	addr = kingpin.Arg("addr", "Server addr").Default(":6000").String()
)

func get_transport_wrapper(name string) TransportWrapper {
	switch name {
	case "none":
		return TTransportWrapper
	case "buffered":
		return NewTBufferedTransportFactory(4096, 4096)
	case "framed":
		return NewTFramedTransportFactory(false, true)
	default:
		panic("invalid transport wrapper")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	kingpin.Parse()

	if *concurrency <= 0 {
		panic("Invalid number of concurrency")
	}

	if *requests <= 0 {
		panic("Invalid number of requests")
	}

	var processor = &Processor{
		pf:        NewTBinaryProtocolFactory(true, true),
		tf:        NewTSocketFactory(*addr),
		tw:        get_transport_wrapper(*transport_wrapper),
		fn:        call("ping"),
		chSuccess: make(chan int, *concurrency*2),
		chError:   make(chan int32, *concurrency*2),
		service:   *service,
	}

	var pipe = make(chan int, *concurrency)
	// collect success messages.
	var pipe1 = make(chan string, 50)
	// collect failed messages.
	var pipe2 = make(chan string, 50)

	go processor.collectSuccess(pipe1)
	go processor.collectError(pipe2)

	fmt.Println("This is Tenchmark, Version 0.1")
	fmt.Println("Copyright 2017 Terence Fan, Baixing, https://github.com/baixing")
	fmt.Println("Licensed under the MIT\n")

	fmt.Printf("Benchmarking %v (be patient)......\n", *addr)

	for i := 0; i < *concurrency; i++ {
		go processor.process(i, pipe)
		wg.Add(1)
	}

	for i := 0; i < *requests; {
		select {
		case pipe <- i:
			i++
		case line := <-pipe1:
			fmt.Println(line)
		}
	}
	close(pipe)
	wg.Wait()

	close(processor.chSuccess)
	close(processor.chError)

	for line := range pipe1 {
		fmt.Println(line)
	}

	for line := range pipe2 {
		fmt.Println(line)
	}
}
