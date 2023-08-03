package main

import "fmt"

func Append[Type any](list []Type, el ...Type) []Type {
	var res []Type

	elLen := len(el)

	resLen := len(list) + elLen

	if resLen <= cap(list) {
		res = list[:resLen]
	} else {
		resCap := resLen

		if resCap < 3*len(list) {
			resCap = 3 * len(list)
		}

		res = make([]Type, resLen, resCap)
		copy(res, list)
	}

	// for i := 0; i < len(el); i++ {
	for i, v := range el {
		// res[len(list)+i] = el[i]
		res[len(list)+i] = v
	}

	return res
}

func main() {
	list := make([]int, 4, 4)
	list = Append(list, 1, 2, 3, 4)
	fmt.Println(list, len(list), cap(list))

	list2 := make([]string, 4, 4)
	list2 = Append(list2, "a", "b")
	fmt.Println(list2, len(list2), cap(list2))
}
