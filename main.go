package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	input := make(chan int)
	go read(input)

	filterNegativeChanel := make(chan int)
	go deleteNegative(input, filterNegativeChanel)
	divTreeChan := make(chan int)
	go removeDivTree(filterNegativeChanel, divTreeChan)

	size := 5
	r := NewRingBuffer(size)
	go writeToBuffer(divTreeChan, r)

	delay := 5
	ticker := time.NewTicker(time.Second * time.Duration(delay))
	go writeToConsole(r, ticker)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal\n", sig)
		os.Exit(0)

	}
}

type RingBuffer struct {
	array []int
	pos   int
	size  int
	m     sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{make([]int, size), -1, size, sync.Mutex{}}
}

func (r *RingBuffer) Push(el int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}
}

func (r *RingBuffer) Get() []int {
	if r.pos <= 0 {
		return nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	var output []int = r.array[:r.pos]
	r.pos = 0
	return output

}

func read(input chan int) {
	for {
		var v int
		_, err := fmt.Scanf("%d\n", &v)
		if err != nil {
			fmt.Println("this is not a number")
		} else {
			input <- v
		}
	}
}

func deleteNegative(currentChan <-chan int, nextChan chan<- int) {
	for number := range currentChan {
		if number >= 0 {
			nextChan <- number
		}
	}
}

func removeDivTree(currentChan <-chan int, nextChan chan<- int) {
	for number := range currentChan {
		if number%3 != 0 {
			nextChan <- number
		}
	}
}
func writeToBuffer(currentChan <-chan int, r *RingBuffer) {
	for number := range currentChan {
		r.Push(number)
	}
}

func writeToConsole(r *RingBuffer, t *time.Ticker) {
	for range t.C {
		buffer := r.Get()
		if len(buffer) > 0 {
			fmt.Println("The Buffer is ", buffer)
		}
	}
}
