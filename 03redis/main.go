package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "118.25.2.25:6379")
	if err != nil {
		fmt.Println("conn redis failed, ", err)
		return
	}
	fmt.Println("redis conn success")
	defer c.Close()

	_, err = c.Do("AUTH", "xu199516..")
	if err != nil {
		fmt.Println("AUTH failed, ", err)
		return
	}

	_, err = c.Do("Set", "abc", 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := redis.Int(c.Do("Get", "abc"))
	if err != nil {
		fmt.Println("get abc failed, ", err)
		return
	}

	r, err = redis.Int(c.Do("Ttl", "abc"))
	fmt.Println(r, err)

	fmt.Println(r)
}
