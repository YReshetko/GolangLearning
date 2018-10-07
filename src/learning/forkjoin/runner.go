package main

import (
	"fmt"
	"learning/forkjoin/task"
	"math/rand"
	"sort"
	"time"
)

func main() {
	fmt.Println("Test regular sort:")
	testPerf(testClasicSorter)
	fmt.Println("Test fork-join sort:")
	testPerf(testFjkSorter)
}

func testPerf(fn func(arr []int)) {
	startTime := time.Now()
	fn(getRandomArr(1000000))
	fmt.Println("Perf time:", time.Now().Sub(startTime))
}

func testFjkSorter(arr []int) {
	task.Process(&sortArray{arr})
}
func testClasicSorter(arr []int) {
	sort.Sort(sortArray{arr})
}

func getRandomArr(count int) []int {
	out := make([]int, count)
	for index, _ := range out {
		out[index] = rand.Intn(1000)
	}
	return out
}

type sortArray struct {
	arr []int
}

func (a sortArray) Len() int           { return len(a.arr) }
func (a sortArray) Swap(i, j int)      { a.arr[i], a.arr[j] = a.arr[j], a.arr[i] }
func (a sortArray) Less(i, j int) bool { return a.arr[i] < a.arr[j] }

func (tsk *sortArray) Execute() {
	if len(tsk.arr) == 2 && tsk.arr[0] > tsk.arr[1] {
		tsk.arr[0], tsk.arr[1] = tsk.arr[1], tsk.arr[0]
	}
}

func (tsk *sortArray) Join(subTasks []task.FjkTask) {
	newArr := make([]int, len(tsk.arr))
	ind1 := 0
	ind2 := 0
	sbTsk1 := subTasks[0].(*sortArray)
	sbTsk2 := subTasks[1].(*sortArray)
	for ind, _ := range newArr {
		if sbTsk1.arr[ind1] > sbTsk2.arr[ind2] {
			newArr[ind] = sbTsk2.arr[ind2]
			ind2++
		} else {
			newArr[ind] = sbTsk1.arr[ind1]
			ind1++
		}
		if ind1 == len(sbTsk1.arr) {
			copy(newArr[ind+1:len(newArr)], sbTsk2.arr[ind2:len(sbTsk2.arr)])
			break
		}
		if ind2 == len(sbTsk2.arr) {
			copy(newArr[ind+1:len(newArr)], sbTsk1.arr[ind1:len(sbTsk1.arr)])
			break
		}
	}
	tsk.arr = newArr
	//fmt.Printf("Merge: %v and %v -> %v\n", sbTsk1.arr, sbTsk2.arr, tsk.arr)
}
func (tsk *sortArray) HasSubtask() bool {
	return len(tsk.arr) > 2
}
func (tsk *sortArray) Fork() []task.FjkTask {
	sbTask1 := &sortArray{tsk.arr[:len(tsk.arr)/2]}
	sbTask2 := &sortArray{tsk.arr[len(tsk.arr)/2:]}
	//fmt.Printf("Create sub tasks: \n%v\n%v\n", sbTask1.arr, sbTask2.arr)
	return []task.FjkTask{sbTask1, sbTask2}
}
