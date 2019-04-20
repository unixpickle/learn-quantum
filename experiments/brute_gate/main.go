package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

func main() {
	results := make(chan quantum.Circuit, 1)
	go Search(3, 15, Toffoli, results)
	for result := range results {
		fmt.Println(result)
	}
}

func Toffoli(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	v := res[0] && res[1]
	if v {
		res[2] = !res[2]
	}
	return res
}
