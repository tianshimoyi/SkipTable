package SkipModel

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//TODO 跳跃链表的查找、插入、删除等操作的期望时间复杂度为O(logn)，效率比链表提高了很多。 跳表的基本性质

//有很多层结构组成
//每一层都是一个有序的链表
//最底层(level 0)的链表包含所有元素
//如果一个元素出现在level i的链表中，则在level i之下的链表都会出现

//TODO 定义节点结构体

type SkipListNode struct {
	key       int             //排序使用的字段
	data      interface{}     //存储的数据
	nextNodes []*SkipListNode //节点指针切片
}

//TODO 定义跳表结构体
type SkipList struct {
	head   *SkipListNode //头节点
	tail   *SkipListNode //尾节点
	length int           //数据总量
	level  int           //层数
	mut    *sync.RWMutex //读写互斥锁
	rand   *rand.Rand    //随机数生成器，用于生成随机层数，随机生成的层数要满足P=0.5的几何分布
}

//TODO 随机生成层数

func (list *SkipList) randomLevel() int {
	level := 1
	for ; level < list.level && list.rand.Uint32()&0x1 == 1; level++ {
	}
	return level
}

//TODO 创建跳表

func NewSkipList(level int) *SkipList {
	if level <= 0 {
		level = 32
	}
	list := &SkipList{
		head:  &SkipListNode{nextNodes: make([]*SkipListNode, level, level)},
		tail:  &SkipListNode{},
		level: level,
		mut:   &sync.RWMutex{},
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
	for index := range list.head.nextNodes {
		list.head.nextNodes[index] = list.tail
	}
	return list
}

//TODO 插入数据

func (list *SkipList) Add(key int, data interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()
	//确定level
	level := list.randomLevel()
	//查找插入部位（从level层开始找）
	update := make([]*SkipListNode, level, level)
	node := list.head
	for index := level - 1; index >= 0; index-- {
		for {
			node1 := node.nextNodes[index]
			if node1 == list.tail || node1.key > key {
				update[index] = node
				break
			} else if node1.key == key {
				node1.data = data
				return
			} else {
				node = node1
			}
		}
	}
	//插入
	sum := 0
	newNode := &SkipListNode{key, data, make([]*SkipListNode, level, level)}
	for index, node := range update {
		node.nextNodes[index], newNode.nextNodes[index] = newNode, node.nextNodes[index]
		sum++
	}
	fmt.Println("sum: ", sum)
	list.length++
}

//TODO 删除数据

func (list *SkipList) Remove(key int) bool {
	list.mut.Lock()
	defer list.mut.Unlock()
	//1. 查找要删除的节点
	node := list.head
	remove := make([]*SkipListNode, list.level, list.level)
	var target *SkipListNode
	for index := len(node.nextNodes) - 1; index >= 0; index-- {
		for {
			node1 := node.nextNodes[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				remove[index] = node
				target = node1
				break
			} else {
				node = node1
			}
		}
	}
	//删除
	if target != nil {
		for index, node1 := range remove {
			if node1 != nil {
				node1.nextNodes[index] = target.nextNodes[index]
			}
		}
		list.length--
		return true
	}
	return false
}

//TODO 查找

func (list *SkipList) Find(key int) interface{} {
	list.mut.Lock()
	defer list.mut.Unlock()
	node := list.head
	for index := list.level - 1; index >= 0; index++ {
		for {
			node1 := node.nextNodes[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == node.key {
				return node1.data
			} else {
				node = node1
			}
		}
	}
	return nil
}

//TODO 遍历跳表

func (list *SkipList) ShowList() {
	list.mut.Lock()
	defer list.mut.Unlock()
	
	for index:=list.level-1;index>=0;index--{
		node:=list.head
		for{
		node1:=node.nextNodes[index]
		if node1==list.tail{
			break
		}else{
			fmt.Printf("%+v\t",node1)
			node=node1
		}
	}
	fmt.Println()
	}
}


func (list *SkipList) Length() int {

	list.mut.RLock()
	defer list.mut.RUnlock()
	return list.length
}