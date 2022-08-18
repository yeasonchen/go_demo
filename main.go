package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i, v := range nums {
		if i2, ok := m[target-v]; ok {
			return []int{i, i2}
		}
		m[v] = i
	}
	return []int{}
}

func main() {
	fmt.Println("haha")
}
