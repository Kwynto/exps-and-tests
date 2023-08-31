package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	qry := "SELECT * FROM entities WHERE id = ? and name = ? LIMIT 1;"

	fmt.Printf("Original: %s \n\n", qry)

	inputData := []any{12, "Kwynto"}

	q1 := strings.Split(qry, "?")
	// fmt.Printf("q1: %v len: %v\n", q1, len(q1))

	lq1 := len(q1)
	if len(inputData) == (lq1 - 1) {

		var q2 string = ""
		if lq1 > 1 {
			for i := 0; i < lq1; i++ {
				if i == (lq1 - 1) {
					q2 = fmt.Sprintf("%v%v", q2, q1[i])
				} else {
					q2 = fmt.Sprintf("%v%v%v", q2, q1[i], inputData[i])
				}

			}
		}
		fmt.Printf("%v\n", q2)
		qry = q2
	} else {
		fmt.Println("not equal placeholders and parameters")
	}

	fmt.Printf("\n---- ---- ---- ---- ---- \n---- ---- ---- ---- ---- \n\n")

	endexp := regexp.MustCompile(`(.+);`)

	if endexp.MatchString(qry) {
		// var q2 string = ""
		// for i, match := range endexp.FindStringSubmatch(qry) {
		// 	if i == 1 {
		// 		q2 = match
		// 	}
		// }
		// matchs := endexp.FindStringSubmatch(qry)
		q2 := endexp.FindStringSubmatch(qry)[1]
		fmt.Printf("%v\n", q2)
		qry = q2
	}

	fmt.Printf("\n---- ---- ---- ---- ---- \n---- ---- ---- ---- ---- \n\n")

	exp01 := regexp.MustCompile(`(?m)^[sS][eE][lL][eE][cC][tT](.+)[fF][rR][oO][mM](.+)[wW][hH][eE][rR][eE](.+)[lL][iI][mM][iI][tT](.+)`)

	fmt.Println(exp01.MatchString(qry))

	for i, match := range exp01.FindStringSubmatch(qry) {
		match = strings.TrimSpace(match)
		fmt.Println(match, "found at index", i)
	}

}
