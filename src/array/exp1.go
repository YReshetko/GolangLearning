package array

import "fmt"

var myArr = []int{1, 2, 6, 1, 7, 9}

func Experiment1() {
	for i, v := range myArr {
		fmt.Printf("myArr[%d] = %d\n", i, v)
	}

	pow := make([]int, 10)

	//Skip value
	for i := range pow {
		//pow[i] = i << uint(i)
		//Equivalent
		pow[i] = i * 2 * 2
	}

	//Skip index
	for _, value := range pow {
		fmt.Printf("Next value = %d\n", value)
	}
	fmt.Printf("Hello")
}
