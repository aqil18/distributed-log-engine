package main

import (
	"fmt"
	"sync"
)

func generateNumbers(total int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // decreases the wg count

	for idx := 1; idx <= total; idx++ {
		fmt.Printf("sending %d to channel\n", idx)
		ch <- idx // <- puts data on the channel
	}
}

func printNumbers(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // decreases the wg count

	for num := range ch {
		fmt.Printf("read %d from channel\n", num)
	}
}

func main() {
	var wg sync.WaitGroup
	numberChan := make(chan int)

	wg.Add(2) // increases count of functions to wait for
	go printNumbers(numberChan, &wg)

	generateNumbers(3, numberChan, &wg)

	close(numberChan)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait() // waits until count is now 0
	fmt.Println("Done!")
}
