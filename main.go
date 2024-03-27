package main

import (
	"fmt"
)

func main() {

	nums := []int{
		1657,
		2030,
		1068,
		753,
		812,
		771,
		1050,
		2287,
		2024,
		2469,
	}

	sum := 0
	for _, num := range nums {
		sum += num
	}

	fmt.Println(sum)
}
