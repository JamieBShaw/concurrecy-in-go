package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var _ io.Reader = &os.File{}

func readCSV(filename string) (<-chan []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	ch := make(chan []string)

	go func() {
		cr := csv.NewReader(file)
		cr.FieldsPerRecord = 3
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

func sanitize(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			if len(val[0]) > 3 {
				continue
			}
			ch <- val
		}
		close(ch)
	}()
	return ch
}

func titleCase(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			val[0] = strings.Title(val[0])
			ch <- val
		}

		close(ch)
	}()

	return ch
}
