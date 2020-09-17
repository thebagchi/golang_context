package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func SlowFunction(caller string, channel chan int) {
	defer func() {
		fmt.Println(caller, "SlowFunction() completed ...")
	}()
	milliseconds := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000)
	fmt.Println("Sleeping for", milliseconds, "ms")
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	if channel != nil {
		channel <- milliseconds
	}
}

func WorkContext(ctx context.Context, channel chan bool) {

	defer func() {
		fmt.Println("WorkContext() completed ...")
		channel <- true
	}()

	sleeper := make(chan int)

	go SlowFunction("WorkContext", sleeper)

	select {
	case <-ctx.Done():
		fmt.Println("SlowFunction() timed out returning")
	case milliseconds := <-sleeper:
		fmt.Println("SlowFunction() returned  normally after", milliseconds, "ms")
	}
}

func DoWork(ctx context.Context) {

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer func() {
		fmt.Println("DoWork() completed ...")
		cancel()
	}()
	channel := make(chan bool)
	go WorkContext(ctx, channel)
	select {
	case <-ctx.Done():
		fmt.Println("WorkContext() timed out returning")
	case <-channel:
		fmt.Println("WorkContext() returned normally")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		fmt.Println("Main() completed ...")
		cancel()
	}()
	DoWork(ctx)
	<-time.After(1 * time.Second)
}
