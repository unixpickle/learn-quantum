package main

import "fmt"

func main() {
	// fmt.Println(Search(3, AllGates(3, false), Toffoli))
	// fmt.Println(Search(3, AllGates(3, true), Or))
	fmt.Println(SearchSqrt(2, AllGates(2, false), CNot))
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
