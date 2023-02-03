package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

func fanOutMain() {
	ch, err := read("combo.csv")
	if err != nil {
		panic(err)
	}

	br1 := fanOut("1", ch)
	br2 := fanOut("2", ch)
	br3 := fanOut("3", ch)

	for {
		if br1 == nil && br2 == nil && br3 == nil {
			break
		}

		select {
		case _, ok := <-br1:
			if !ok {
				br1 = nil
			}
		case _, ok := <-br2:
			if !ok {
				br2 = nil
			}
		case _, ok := <-br3:
			if !ok {
				br3 = nil
			}
		}
	}

	fmt.Println("\nAll complete")
}

func fanOut(worker string, ch <-chan []string) chan struct{} {
	chE := make(chan struct{})

	read := func(ch <-chan []string) {
		for v := range ch {
			fmt.Println(worker, v)
		}
	}

	go func() {
		read(ch)
		close(chE)
	}()

	return chE
}

func read(filename string) (<-chan []string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file %w", err)
	}

	ch := make(chan []string)
	cr := csv.NewReader(f)

	go func() {
		for {
			record, err := cr.Read()
			if errors.Is(err, io.EOF) {
				close(ch)
				return
			}

			ch <- record
		}
	}()

	return ch, nil
}
