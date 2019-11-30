package main

import (
	"fmt"
	"github.com/robfig/cron"
)

func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func main() {
	i := 0
	c := newWithSeconds()
	spec := "*/1 * * * * *"
	c.AddFunc(spec, func() {

		i++
		fmt.Println("cron running: ", i)
	})
	c.Start()
	select {}

}
