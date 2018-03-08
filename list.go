package main

import (
	"container/list"
	"flag"
	"fmt"
	"time"
	"unsafe"
)

type kernelList struct {
	prev, next *kernelList
	length uint64
}

func NewList() *kernelList {
	var l = &kernelList{}
	l.prev = l
	l.next = l
	l.length = 0
	return l
}

func (l *kernelList) add(cur, prev, next *kernelList) {
	prev.next = cur
	next.prev = cur
	cur.prev = prev
	cur.next = next
	l.length++
}

func (l *kernelList) Len() uint64 {
	return l.length
}

func (l *kernelList) PushFront(cur *kernelList) {
	l.add(cur, l, l.next)
}

func (l *kernelList) PushBack(cur *kernelList) {
	l.add(cur, l.prev, l)
}

func (l *kernelList) InsertAfter(cur, pend *kernelList) {
	l.add(pend, cur, cur.next)
}

func (l *kernelList) InsertBefore(pend, cur *kernelList) {
	l.add(pend, cur.prev, cur)
}

func (l *kernelList) Remove(cur *kernelList) {
	l.length--
	cur.prev.next = cur.next
	cur.next.prev = cur.prev
	cur.prev = nil
	cur.next = nil
}

func offsetOf() uintptr {
	var val = &myData{}
	return unsafe.Offsetof(val.kernelList)
}

func listEntry(ptr *kernelList) unsafe.Pointer {
	//fmt.Printf("ptr-addr:%v, ptr-uintptr:%v, offsetof:%v\n", unsafe.Pointer(ptr), uintptr(unsafe.Pointer(ptr)), offsetOf())
	return unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - offsetOf())
}

//just for test
type stdData struct {
	id   int
	age  int
	name string
}

//std list
func stdList(size int) *list.List {
	var start = time.Now()
	var head = list.New()
	for i := 0; i < size; i++ {
		var sd = stdData{
			id:   i,
			name: fmt.Sprintf("name_%d", i),
			age:  i + 100,
		}
		head.PushFront(sd)
	}

	for i, n := size, 2*size; i < n; i++ {
		var sd = stdData{
			id:   i,
			name: fmt.Sprintf("name_%d", i),
			age:  i + 100,
		}
		head.PushBack(sd)
	}
	fmt.Printf("std-list task time:%.3dms\n", time.Now().Sub(start).Nanoseconds()/1000.0/1000.0)

	return head
}

func stdForEach(head *list.List) {
	for pos := head.Front(); pos != nil; pos = pos.Next() {
		fmt.Printf("pos:%+v, entry:%+v\n", pos, pos.Value)
	}
}

type myData struct {
	id   int
	age  int
	name string
	kernelList
}

//kernel list
func kerList(size int) *kernelList {
	var start = time.Now()
	var head = NewList()
	for i := 0; i < size; i++ {
		var md = myData{
			id:   i,
			name: fmt.Sprintf("name_%d", i),
			age:  i + 100,
		}
		head.PushFront(&md.kernelList)
	}

	for i, n := size, 2*size; i < n; i++ {
		var md = myData{
			id:   i,
			name: fmt.Sprintf("name_%d", i),
			age:  i + 100,
		}
		head.PushBack(&md.kernelList)
	}
	fmt.Printf("kernel-list task time:%.3dms\n", time.Now().Sub(start).Nanoseconds()/1000.0/1000.0)

	return head
}

func kerForEach(head *kernelList) {
	var entry *myData
	for pos := head.next; pos != head; pos = pos.next {
		entry = (*myData)(listEntry(pos)) //list在myData的任意字段位置
		fmt.Printf("pos:%+v, entry:%+v\n", *pos, *entry)
	}
}

func multiOpsList(size int) {
	var head = NewList()

	var one = myData{
		id:   1,
		name: fmt.Sprintf("name_%d", 1),
		age:  1 + 100,
	}
	var two = myData{
		id:   2,
		name: fmt.Sprintf("name_%d", 2),
		age:  2 + 100,
	}
	var three = myData{
		id:   3,
		name: fmt.Sprintf("name_%d", 3),
		age:  3 + 100,
	}
	var four = myData{
		id:   4,
		name: fmt.Sprintf("name_%d", 4),
		age:  4 + 100,
	}
	var five = myData{
		id:   5,
		name: fmt.Sprintf("name_%d", 5),
		age:  5 + 100,
	}
	head.PushFront(&two.kernelList)
	kerForEach(head)
	fmt.Printf("list size:%d\n", head.Len())
	println()

	head.InsertBefore(&one.kernelList, &two.kernelList)
	kerForEach(head)
	fmt.Printf("list size:%d\n", head.Len())
	println()

	head.InsertAfter(&two.kernelList, &three.kernelList)
	kerForEach(head)
	fmt.Printf("list size:%d\n", head.Len())
	println()

	head.PushBack(&four.kernelList)
	kerForEach(head)
	fmt.Printf("list size:%d\n", head.Len())
	println()

	head.InsertAfter(&four.kernelList, &five.kernelList)
	kerForEach(head)
	fmt.Printf("list size:%d\n", head.Len())
}

func main() {
	var xrange = flag.Int("r", 0, "link size range")
	flag.Parse()
	fmt.Printf("create link size:%d\n\n", *xrange)

	if false {
		var sh = stdList(*xrange)
		stdForEach(sh)
	}

	if true {
		var kh = kerList(*xrange)
		kerForEach(kh)
	}

	if false {
		multiOpsList(*xrange)
	}
	return
}
