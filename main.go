package main

import (
	"fmt"
	"runtime"
)

type MetaCube struct {
	i    int
	cube int32
}

const numCount = 1000

func main() {
	allNums := make([]MetaCube, 0, numCount)
	{
		var i, j int
		for i = 0; i < numCount; i++ {
			j = i + 1
			allNums = append(allNums, MetaCube{j, fx(j)})
		}
	}

	receiver := make(chan [4]int, runtime.GOMAXPROCS(0))

	go func() {
		for a := 1; a < numCount; a++ {
			for b := a + 1; b < numCount-1; b++ {
				if b < a {
					continue
				}
				go find(allNums[a], allNums[b], allNums[a+1:b], receiver)
			}
		}
		close(receiver)
	}()

	var iterId int
	for out := range receiver {
		iterId++
		fmt.Printf("%d)\t%d^3 + %d^3 = %d^3 + %d^3\n", iterId, out[0], out[1], out[2], out[3])
	}
	fmt.Println("the end")
}

//find поиск корней уравнения
func find(a, b MetaCube, nums []MetaCube, receiver chan [4]int) {
	sum := a.cube + b.cube
	for len(nums) > 1 {
		c := nums[0]
		nums = nums[1:]
		for _, d := range nums {
			if c.cube+d.cube == sum {
				receiver <- [4]int{a.i, b.i, c.i, d.i}
			}
		}
	}
}

func fx(i int) int32 {
	return int32(i * i * i)
}
