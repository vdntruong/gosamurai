package subtleties

import "fmt"

func RangeOverInteger(end int) {
	for i := range end {
		fmt.Println(i + 1)
	}
}

/*
1
2
3
4
5
6
7
8
9
10
*/
