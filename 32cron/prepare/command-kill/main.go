package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		err        error
		resultChan chan *result
		res        *result
	)

	resultChan = make(chan *result, 1000)

	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
		)

		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "ls;sleep 2;ls -al")
		output, err = cmd.CombinedOutput()

		resultChan <- &result{
			err:   err,
			ouput: string(output),
		}

	}()

	time.Sleep(1 * time.Second)
	cancelFunc()

	res = <-resultChan
	fmt.Println("ok", res.ouput,res.err)
}

type result struct {
	err   error
	ouput string
}
