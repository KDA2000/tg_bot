package main

import (
	"fmt"
	"time"
)

func main() {
	readere(doubler(writerr()))
}

func writerr() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i + 1
		}
		close(ch)
	}()
	return ch
}

func doubler(value <-chan int) chan int {
	ch := make(chan int)
	go func() {
		for j := range value {
			time.Sleep(500 * time.Millisecond)
			ch <- j * 2
		}
		close(ch)
	}()
	return ch
}

func readere(read <-chan int) {
	for v := range read {
		fmt.Println(v)
	}
}
