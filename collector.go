package main

import (
	"fmt"
	"math"
	"time"

	. "github.com/stdrickforce/thriftgo/thrift"
)

func collect_success(ch <-chan int, pipe chan<- string) {
	defer close(pipe)

	snano := time.Now().UnixNano()

	var (
		s     = make([]int, 0)
		count = 0
	)

	for duration := range ch {
		count++
		if count%1000 == 0 {
			pipe <- fmt.Sprintf("Completed %d requests", count)
		}
		s = append(s, duration)
	}
	pipe <- fmt.Sprintf("Finished %d requests", count)
	pipe <- ""

	dnano := time.Now().UnixNano() - snano

	l := len(s)
	sort(s, 0, l-1)

	v := func(denominator int) float64 {
		if denominator <= 0 {
			return float64(s[l-1]) / 1000
		} else {
			return float64(s[l*(denominator-1)/denominator]) / 1000
		}
	}

	var (
		duration = float64(dnano) / math.Pow(10, 9)
		qps      = float64(l) / duration
	)

	pipe <- fmt.Sprintf("%-24s%s", "Server Address:", *addr)
	pipe <- ""
	pipe <- fmt.Sprintf("%-24s%d", "Concurrency level:", *concurrency)
	pipe <- fmt.Sprintf("%-24s%.3f seconds", "Time taken for tests:", duration)
	pipe <- fmt.Sprintf("%-24s%d", "Complete requests:", l)
	pipe <- fmt.Sprintf("%-24s%d", "Failed requests:", *requests-l)
	pipe <- fmt.Sprintf("%-24s%.2f [#/sec] (mean)", "Request per second:", qps)
	pipe <- ""

	if l == 0 {
		return
	}

	pipe <- "Percentage of the requests served within a certain time (ms)"
	pipe <- fmt.Sprintf("%4d%% %8.2f", 50, v(2))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 66, v(3))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 75, v(4))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 80, v(5))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 90, v(10))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 95, v(20))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 98, v(50))
	pipe <- fmt.Sprintf("%4d%% %8.2f", 99, v(100))
	pipe <- fmt.Sprintf("%4d%% %8.2f (longest request)", 100, v(-1))
	pipe <- ""
}

func collect_error(ch <-chan int32, pipe chan<- string) {
	defer close(pipe)

	var (
		count        int
		distribution = make(map[int32]int)
	)

	for mtype := range ch {
		count++
		distribution[mtype]++
	}
	var s = func(k int32) string {
		switch k {
		case ExceptionUnknown:
			return "ExceptionUnknown"
		case ExceptionUnknownMethod:
			return "ExceptionUnknownMethod"
		case ExceptionInvalidMessageType:
			return "ExceptionInvalidMessageType"
		case ExceptionWrongMethodName:
			return "ExceptionWrongMethodName"
		case ExceptionBadSequenceID:
			return "ExceptionBadSequenceID"
		case ExceptionMissingResult:
			return "ExceptionMissingResult"
		case ExceptionInternalError:
			return "ExceptionInternalError"
		case ExceptionProtocolError:
			return "ExceptionProtocolError"
		default:
			return fmt.Sprintf("%d", k)
		}
	}

	if count > 0 {
		pipe <- fmt.Sprintf("Count of the exception replied by server:")
		for mtype, val := range distribution {
			pipe <- fmt.Sprintf("%-32s%d", s(mtype), val)
		}
	}
}

func collect(ch1 <-chan int, ch2 <-chan int32) {
	wg2.Add(1)
	defer wg2.Done()

	pipe1, pipe2 := make(chan string, 50), make(chan string, 50)
	go collect_success(ch1, pipe1)
	go collect_error(ch2, pipe2)

	for line := range pipe1 {
		fmt.Println(line)
	}

	for line := range pipe2 {
		fmt.Println(line)
	}
}
