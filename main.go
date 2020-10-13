package main

import (
	m "DataStructure/SkipTable/SkipModel"
	"fmt"
)

func main() {
	list := m.NewSkipList(5)
	list.Add(3, []int{10, 9, 5, 4, 3})
	list.Add(9, []int{10, 9, 5, 4, 3})
	list.Add(1, []int{10, 9, 5, 4, 3})
	list.Add(7, []int{10, 9, 5, 4, 3})
	fmt.Println(list.Length())
	list.ShowList()
}
