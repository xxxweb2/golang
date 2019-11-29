package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

func Counter(wg *sync.WaitGroup) {
	time.Sleep(time.Second)

	var counter int
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 200)
		counter++
		fmt.Println(counter)
	}
	wg.Done()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {

	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
