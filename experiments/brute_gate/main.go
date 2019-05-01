package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

func main() {
	fmt.Println(Search(3, AllGates(3, false), quantum.NewClassicalGate(Toffoli, "")))
	// fmt.Println(Search(3, AllGates(3, true), Or))
	// fmt.Println(SearchSqrt(2, AllGates(2, false), CNot))
	// fmt.Println(SearchSqrt(2, AllGates(2, false), Swap))
	// fmt.Println(Search(3, AllGates(3, false), CSwap))
	// fmt.Println(SearchSqrt(1, AllGates(1, false), Not))
	// fmt.Println(SearchCtrl(2, AllGates(2, false), &quantum.SqrtNotGate{Bit: 1}))
}

func Not(b []bool) []bool {
	res := make([]bool, 1)
	copy(res, b)
	res[0] = !res[0]
	return res
}

func Toffoli(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	if res[0] && res[1] {
		res[2] = !res[2]
	}
	return res
}

func Or(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	if res[0] || res[1] {
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

func Swap(b []bool) []bool {
	res := make([]bool, 2)
	res[0] = b[1]
	res[1] = b[0]
	return res
}

func CSwap(b []bool) []bool {
	res := make([]bool, 3)
	copy(res, b)
	if b[0] {
		res[1] = b[2]
		res[2] = b[1]
	}
	return res
}
