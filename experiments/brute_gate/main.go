package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

func main() {
	results := make(chan quantum.Circuit, 1)
	go Search(3, 20, Toffoli, results)
	// go SearchSqrt(2, 15, CNot, results)
	for result := range results {
		fmt.Println(result)
	}
}

func Toffoli(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	if res[0] && res[1] {
		res[2] = !res[2]
	}
	return res
}

func CNot(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	if res[0] {
		res[1] = !res[1]
	}
	return res
}
