package main

import (
	"fmt"
	"sync"
	"time"
)

func broker(input chan int, numOfWorkers int) chan int {
	if input == nil {
		return nil
	}
	var wg sync.WaitGroup
	output := make(chan int)
	wg.Add(numOfWorkers)

	for i := 0; i < numOfWorkers; i++ {
		go func(workerId int) {
			defer wg.Done()
			for {
				val, ok := <-input
				fmt.Printf("worker %d received value %d\n", workerId, val)
				if !ok {
					fmt.Printf("worker %d stopped\n", workerId)
					return
				}
				time.Sleep(time.Millisecond * 100)
				output <- val * val
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}


func main() {
	const N int = 20
	const NUM_OF_FANIN = 3

	input := make(chan int, N)

	go func() {
		for i := 0; i < N; i++ {
			input <- i
		}
		close(input)
	}()

	go func() {
		for i := 0; i < NUM_OF_FANIN; i++ {
			go func(workerId int) {
				for {

				}
			}(i)
		}
	}()

	output := broker(input, 3)

	out := make([]int, 0)
	for val := range output {
		out = append(out, val)
	}

	for _, val := range out {
		fmt.Printf("output: %d\n", val)
	}
}
