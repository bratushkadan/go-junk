package main

package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	ch := make(chan uint64)

	const n int = 30_500_000

	start := time.Now()
	go SumOfSquares(ch, n)

	fmt.Printf("For n=%d, sum=%d\n", n, <-ch)
	fmt.Printf("%dmu elapsed.\n", time.Now().Sub(start).Microseconds())
}

func sumSquaresRange(ch chan<- uint64, from, to int) {
	var s uint64
	var i int = from

	for ; i < to; i++ {
		s += uint64(i) * uint64(i)
	}

	ch <- s
}

func SumOfSquares(ch chan<- uint64, n int) {
	cpuNum := runtime.NumCPU()

	internalCh := make(chan uint64, cpuNum)

	chunkSize := n / int(cpuNum)

	var i int
	for ; i < cpuNum; i++ {
		go sumSquaresRange(internalCh, i*chunkSize+1, (i+1)*chunkSize)
	}

	var s uint64
	finished := 0
	for finished != cpuNum {
		s += <-internalCh
		finished++
	}

	ch <- s
}

/*
	* single-threaded *
	For n=30500000, sum=12808235802381322608
	17975mu elapsed.

	* multi-threaded *
	For n=30500000, sum=12805584589881322608
	2344mu elapsed.
*/

