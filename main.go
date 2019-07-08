package main

import (
	"fmt"
	"sync"
)

func main() {
	res := fanin(fanout(generateNumbers(), squareNumbers))

	for v := range res {
		fmt.Println(v)
	}

	fmt.Println("finished")
}

type Op func(chan interface{}) chan interface{}

func generateNumbers() chan interface{} {
	fmt.Println("generateNumbers called")
	out := make(chan interface{}, 100)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("gen: ", i)
			out <- i
		}
		close(out)
	}()

	return out
}

func fanout(in chan interface{}, fct Op) []chan interface{} {
	fmt.Println("fanout called")
	var out []chan interface{}

	for i := 0; i < 5; i++ {
		res := fct(in)
		out = append(out, res)
	}

	return out
}

func squareNumbers(in chan interface{}) chan interface{} {
	fmt.Println("squareNumbers called")
	out := make(chan interface{}, 100)

	go func() {
		for v := range in {
			num, _ := v.(int)

			squared := num * num
			out <- squared
		}
		close(out)
	}()

	return out
}

func fanin(ins []chan interface{}) chan interface{} {
	fmt.Println("fanin called")
	out := make(chan interface{}, 100)
	var waitgroup sync.WaitGroup
	waitgroup.Add(len(ins))

	for _, v := range ins {
		go func(w chan interface{}) {
			for val := range w {
				out <- val
			}
			waitgroup.Done()
		}(v)
	}

	go func() {
		waitgroup.Wait()
		close(out)
	}()

	return out
}
