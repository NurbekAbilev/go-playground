package main

import (
	"fmt"
	"sync"
)

func fanout(input <-chan int, size int) []chan int {
	fanout := make([]chan int, 0)
	for i := 0; i < size; i++ {
		ch := make(chan int)
		fanout = append(fanout, ch)
	}

	go func() {
		for val := range input {
			for _, ch := range fanout {
				ch <- val
			}
		}

		for _, ch := range fanout {
			close(ch)
		}
	}()

	return fanout
}

func worker(name string, input <-chan int) {
	for val := range input {
		fmt.Printf("worker [%s] got value %d\n", name, val)
	}
}

func main() {
	input := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			input <- i
		}
		close(input)
	}()

	const SIZE int = 2
	fch := fanout(input, SIZE)

	var wg sync.WaitGroup

	for i, ch := range fch {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(fmt.Sprintf("Worker %d", i), ch)
		}()
	}

	wg.Wait()
}
