package main

import (
	"reflect"
)

func tee(done, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func tee2(done, in <-chan interface{}, exitChannelLen int) []chan interface{} {
	returnChannels := make([]chan interface{}, exitChannelLen)

	teeUp := func(done <-chan interface{}, c chan interface{}) {
		defer close(c)
		for val := range orDone(done, in) {
			var exitChan = c
			select {
			case <-done:
			case exitChan <- val:
			}
		}
	}

	for _, c := range returnChannels {
		go teeUp(done, c)
	}

	return returnChannels
}

type SimpleOutChannel interface {
	Out() <-chan interface{} // The readable end of the channel.
}

type SimpleInChannel interface {
	In() chan<- interface{} // The writeable end of the channel.
	Close()                 // Closes the channel. It is an error to write to In() after calling Close().
}

func tee3(input SimpleOutChannel, outputs []SimpleInChannel, closeWhenDone bool) {
	cases := make([]reflect.SelectCase, len(outputs))
	for i := range cases {
		cases[i].Dir = reflect.SelectSend
	}
	for elem := range input.Out() {
		for i := range cases {
			cases[i].Chan = reflect.ValueOf(outputs[i].In())
			cases[i].Send = reflect.ValueOf(elem)
		}
		for range cases {
			chosen, _, _ := reflect.Select(cases)
			cases[chosen].Chan = reflect.ValueOf(nil)
		}
	}
	if closeWhenDone {
		for i := range outputs {
			outputs[i].Close()
		}
	}
}
