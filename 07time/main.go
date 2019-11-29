package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	fmt.Println(now)
	fmt.Println("year: ", year)
	fmt.Println("month: ", month)
	fmt.Println("day: ", day)
	fmt.Println("hour: ", hour)

	fmt.Println("minute: ", minute)
	fmt.Println("second: ", second)

	fmt.Println("year", now.Format("2006"))
	fmt.Println("month", now.Format("1"))
	fmt.Println("day", now.Format("2"))
	fmt.Println("hour", now.Format("3"))
	fmt.Println("min", now.Format("04"))


}
