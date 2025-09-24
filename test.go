package main

import (
	"fmt"
	"sync"
	"time"
)

func ti() {
	count := 0
	tt := time.Now()
	var wg sync.WaitGroup
	var mu sync.RWMutex

	for i := 0; i < 500; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			time.Sleep(time.Nanosecond)
			mu.RLock()
			_ = count
			mu.RUnlock()
		}()
		go func() {
			defer wg.Done()
			time.Sleep(time.Nanosecond)
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
	fmt.Println(time.Now().Sub(tt).Seconds())
}

type Person struct {
	name string
	age  int
}

func (p Person) greet() {
	fmt.Println("Hello ", p.name)
}

type Greeter interface {
	greet()
}

func sayhello(g Greeter) {
	g.greet()
}

func main() {
	g := Person{name: "asda"}
	sayhello(g)
}

func what() {
	start := time.Now()
	num := 0
	for i := 0; i < 500; i++ {
		time.Sleep(time.Nanosecond)
		num++
	}
	fmt.Println(num)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func gor() {
	start := time.Now()
	num := 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Nanosecond)
			mu.Lock()
			num++
			mu.Unlock()
		}()

	}
	wg.Wait()
	fmt.Println(num)
	fmt.Println(time.Now().Sub(start).Seconds())
}
