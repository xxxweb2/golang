package main

import (
	"fmt"
	"unsafe"
)

type Person struct {
	Name string
	age  int
	sex  int
}

func main() {
	p := Person{
		Name: "许欣欣",
	}

	base := uintptr(unsafe.Pointer(&p))
	offset := unsafe.Offsetof(p.sex)
	ptr := unsafe.Pointer(base + offset)
	*(*int)(ptr) = 3

	fmt.Println(p)
}
