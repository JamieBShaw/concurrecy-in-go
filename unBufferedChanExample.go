package main

import (
	"fmt"
)

func unBuggeredChanExample() {
	ch := make(chan int)
	go worker(ch)

	fmt.Println("Notifying worker to begin")
	ch <- 1

	fmt.Println("Wait for worker to be done")
	<- ch
	fmt.Println("Worker done")
}

func worker(ch chan int) {

	fmt.Println("Waiting for main to tell us to begin")
	<-ch

	fmt.Println("Worker starts working")
	for i := 0; i < 9; i++ {
		fmt.Println("Working hard: ", i)
		if i == 8 {
			fmt.Println("Work Done")
			ch <- 1
		}
	}

	fmt.Println("Wait for main to signal we can restart work")
	<-ch
}
