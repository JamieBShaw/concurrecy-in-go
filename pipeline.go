package main

import "fmt"

func pipeline() {

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
	pipepline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipepline {
		fmt.Println(v)
	}

}

func generator(done <-chan interface{}, numbers ...int) <-chan int {
	intStream := make(chan int, len(numbers))

	go func() {
		defer close(intStream)
		for _, i := range numbers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}

	}()
	return intStream
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)

	go func() {
		defer close(multipliedStream)

		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()

	return multipliedStream
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)

	go func() {
		defer close(addedStream)

		for i := range intStream {
			select {
			case <-done:
				return

			case addedStream <- i + additive:
			}
		}
	}()

	return addedStream
}
