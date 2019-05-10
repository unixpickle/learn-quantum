package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

func main() {
	// fmt.Println(Search(3, AllGates(3, false), quantum.NewClassicalGate(Toffoli, "")))
	// fmt.Println(SearchSqrt(1, AllGates(1, false), quantum.NewClassicalGate(Not, "")))
	// fmt.Println(SearchCtrl(2, AllGates(2, false), &quantum.SqrtNotGate{Bit: 1}))
	// fmt.Println(SearchCtrl(2, AllGates(2, false), &quantum.HGate{Bit: 1}))
	idx := 3
	fmt.Println(Search(4, AllGates(4, true), &constAdder{
		Value:  5,
		Target: quantum.Reg{0, 1, 2},
		Carry:  &idx,
	}))
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
