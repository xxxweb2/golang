package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num float64 = 1.23456
	val := reflect.ValueOf(num)
	fmt.Println("type: ", reflect.TypeOf(num))
	fmt.Println("value: ", val)
}
