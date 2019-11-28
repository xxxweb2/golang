package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", "118.25.2.25:6379")
			if err != nil {
				return nil, err
			}

			_, err = c.Do("AUTH", "xu199516..")

			return c, err
		},
	}
}

func main() {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("Set", "abc", 200)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := redis.Int(c.Do("Get", "abc"))
	if err != nil {
		fmt.Println("get abc failed: ", err)
		return
	}
	fmt.Println(r)
	pool.Close()
}
