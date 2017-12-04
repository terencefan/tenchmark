package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"xparser"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/transport"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	wg  sync.WaitGroup
	wg2 sync.WaitGroup
)

var (
	app = kingpin.New("tenchmark", "Thrift benchmark command-line tools")

	protocol = app.Flag("protocol", "Specify protocol factory").Default("binary").String()

	run   = app.Command("run", "Run benchmark tests")
	build = app.Command("build", "Build cases from thrift file and json inputs")

	// run command args.
	service     = run.Flag("service", "Specify service name (multiplexed)").String()
	requests    = run.Flag("requests", "Number of requests to perform").Short('n').Default("1000").Int()
	concurrency = run.Flag("concurrency", "Number of multiple requests to make at a time").Short('c').Default("10").Int()
	transport   = run.Flag("transport", "Specify transport factory").Default("socket").String()
	wrapper     = run.Flag("wrapper", "Specify transport wrapper").Default("buffered").String()
	casefile    = run.Flag("case", "Generated from `tenchmark build`").Short('b').ExistingFile()
	addr        = run.Arg("addr", "Server addr").Default(":6000").String()

	// build command args
	api_json    = build.Flag("json", "Path to api json file").Required().ExistingFile()
	outputdir   = build.Flag("out", "Path to generated .in files").Default("cases").String()
	multiplexed = build.Flag("multiplexed", "If service protocol is Multiplexed.").Bool()
	thrift_file = build.Arg("thrift", "Path to thrift file").Required().ExistingFile()
)

func get_transport_factory(name, addr string) TransportFactory {
	switch name {
	case "socket":
		return NewTSocketFactory(addr)
	case "unix":
		return NewTUnixSocketFactory(addr)
	case "http":
		return NewTHttpTransportFactory(addr)
	default:
		panic("invalid transport type: " + name)
	}
}

func get_transport_wrapper(name string) TransportWrapper {
	switch name {
	case "buffered":
		return NewTBufferedTransportFactory(4096, 4096)
	case "framed":
		return NewTFramedTransportFactory(false, true)
	default:
		panic("invalid transport wrapper: " + name)
	}
}

func get_protocol_factory(name string) ProtocolFactory {
	switch name {
	case "binary":
		return NewTBinaryProtocolFactory(true, true)
	default:
		panic("invalid protocol: " + name)
	}
}

func run_test() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *concurrency <= 0 {
		panic("Invalid number of concurrency")
	}

	if *requests <= 0 {
		panic("Invalid number of requests")
	}

	var processor = &Processor{
		service:   *service,
		pf:        NewTBinaryProtocolFactory(true, true),
		tf:        get_transport_factory(*transport, *addr),
		tw:        get_transport_wrapper(*wrapper),
		chSuccess: make(chan int, *concurrency*2),
		chError:   make(chan int32, *concurrency*2),
	}

	if err := processor.initMessage(*casefile); err != nil {
		panic(err)
	}

	if err := processor.test(); err != nil {
		panic(err)
	}

	var pipe = make(chan int, *concurrency)

	fmt.Println("This is Tenchmark, Version 0.1")
	fmt.Println("Copyright 2017 Terence Fan, Baixing, https://github.com/baixing")
	fmt.Println("Licensed under the MIT\n")

	fmt.Printf("Benchmarking %v (be patient)......\n", *addr)

	for i := 0; i < *concurrency; i++ {
		go processor.process(i, pipe)
		wg.Add(1)
	}
	go collect(processor.chSuccess, processor.chError)

	for i := 0; i < *requests; i++ {
		pipe <- i
	}
	close(pipe)
	wg.Wait()

	close(processor.chSuccess)
	close(processor.chError)
	wg2.Wait()
}

func build_cases() {
	parser, err := xparser.InitParser(*thrift_file)
	if err != nil {
		panic(err)
	}

	os.Mkdir(*outputdir, os.FileMode(0755))

	apis, err := xparser.NewApiParser(*api_json)
	if err != nil {
		panic(err)
	}

	for name, apicase := range apis.GetCases() {
		filename := fmt.Sprintf("%s/%s.in", *outputdir, name)
		trans := xparser.NewFileOutputStream(filename)
		proto := get_protocol_factory(*protocol).GetProtocol(trans)

		if err = trans.Open(); err != nil {
			panic(err)
		}
		if err = parser.BuildRequest(proto, apicase, *multiplexed); err != nil {
			panic(err)
		}
		if err = trans.Close(); err != nil {
			panic(err)
		}
		fmt.Printf("%s sucessfully generated.\n", filename)
	}

}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		run_test()
	case build.FullCommand():
		build_cases()
	}
}
