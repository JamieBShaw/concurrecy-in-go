package main

import (
	"fmt"
	"sync"
	"time"
)

func fanInMain() {
	file1, err := read("file1.csv")
	if err != nil {
		panic(err)
	}

	file2, err := read("file2.csv")
	if err != nil {
		panic(err)
	}

	merged := merge(file1, file2)

	shutdown := make(chan struct{})

	go func() {
		for v := range merged {
			fmt.Println(v)
		}

		close(shutdown)
	}()

	<-shutdown
}

func merge(cs ... <-chan []string) <-chan []string {
	chans := len(cs)
	merged := make(chan []string)
	wait := make(chan struct{}, chans)

	send := func(c <-chan []string) {
		defer func() {
			wait <- struct{}{}
		}()
		for n := range c {
			merged <- n
		}
	}
	for _, c := range cs {
		go send(c)
	}

	go func() {
		for range wait {
			chans--
			if chans == 0 {
				break
			}
		}
		close(merged)
	}()
	return merged
}


func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{}, len(channels))

	multiplexer := func(c <- chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplexer(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

//func main() {
//	ch := fanIn(generator("Hello"), generator("Bye"))
//	for i := 0; i < 10; i++ {
//		fmt.Println(<- ch)
//	}
//}

// fanIn is itself a generator
func fanInMore(ch1, ch2 <-chan string) <-chan string { // receives two read-only channels
	new_ch := make(chan string)
	go func() { for { new_ch <- <-ch1 } }() // launch two goroutine while loops to continuously pipe to new channel
	go func() { for { new_ch <- <-ch2 } }()
	return new_ch
}

func generatorMore(msg string) <-chan string { // returns receive-only channel
	ch := make(chan string)
	go func() { // anonymous goroutine
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Second)
		}
	}()
	return ch
}
