package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func condExample() {
	lock := sync.Mutex{}
	lock.Lock()
	cond := sync.NewCond(&lock)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()

		fmt.Println("First go routine has started and waits for 1 second before broadcasting condition")

		time.Sleep(1 * time.Second)

		fmt.Println("First go routine broadcasts condition")

		cond.Broadcast()
	}()

	go func() {
		defer waitGroup.Done()

		fmt.Println("Second go routine has started and is waiting on condition")

		cond.Wait()

		fmt.Println("Second go routine unlocked by condition broadcast")
	}()

	fmt.Println("Main go routine starts waiting")

	waitGroup.Wait()

	fmt.Println("Main go routine ends")
}

func differentCondExample() {
	var age = make(map[string]int)

	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	// listener 1
	go listen("lis1", age, cond)

	// listener 2
	go listen("lis2", age, cond)

	// broadcast
	go broadcast(age, cond)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func listen(name string, a map[string]int, c *sync.Cond) {
	c.L.Lock()
	c.Wait()
	fmt.Println(name, " age:", a["T"])
	c.L.Unlock()
}

func broadcast(a map[string]int, c *sync.Cond) {
	time.Sleep(time.Second)
	c.L.Lock()
	a["T"] = 25
	c.Broadcast()
	c.L.Unlock()
}
