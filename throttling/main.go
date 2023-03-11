package main

import (
	"fmt"
	"time"
	"sync"
)

func downloadJSON(u string) {
	println(u)
	time.Sleep(time.Second * 6)
}

func main() {
	before := time.Now()

	limit := make(chan struct{}, 20)

	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			limit <- struct{}{}
			defer wg.Done()

			u := fmt.Sprintf("http://example.com/api/users?id=%d", i)
			downloadJSON(u)

			<-limit
		}()
	}

	wg.Wait()

	fmt.Printf("%v\n", time.Since(before))
}