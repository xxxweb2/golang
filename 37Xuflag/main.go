package main

import (
	"fmt"
	"study/golang/37Xuflag/flag"
)

func main()  {
	var age = xuFlag.Int("age", 18, "age age")
	var height = xuFlag.Int("height", 20, "height height")

	xuFlag.Parse()
	//flag.Usage()
	fmt.Printf("age=%d  height=%d", *age, *height)
}
