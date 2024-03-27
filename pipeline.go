package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Окончательная версия проекта с двухступенчатым докером

func main() {
	input := make(chan int)
	go read(input)
	log.Println("\nStep #1\nEnter 5 numbers into the buffer")

	filterNegativeChanel := make(chan int)
	go deleteNegative(input, filterNegativeChanel)
	divTreeChan := make(chan int)
	go removeDivTree(filterNegativeChanel, divTreeChan)

	size := 6
	r := NewRingBuffer(size)
	go writeToBuffer(divTreeChan, r)

	delay := 6
	ticker := time.NewTicker(time.Second * time.Duration(delay))
	go writeToConsole(r, ticker)
	c := make(chan os.Signal, 6)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	log.Printf("Got %s signal\n", sig)
	os.Exit(0)
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
	output := r.array[:r.pos]
	r.pos = 0
	return output

}

func read(input chan int) {
	for {
		var v int
		_, err := fmt.Scanf("%d\n", &v)
		if err != nil {
			log.Println("\nStep No. 404\nthis is not a number")
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
			log.Println("\nThe Buffer is", buffer)
		}
	}
}
