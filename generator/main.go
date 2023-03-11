package main

import (
	"fmt"
	// "time"
	"math/rand"
)

// func fanInV1(ch1, ch2 <- chan string) <- chan string {
// 	new_ch := make(chan string)
// 	go func ()  {
// 		for {
// 			new_ch <- <-ch1
// 		}
// 	}()
// 	go func ()  {
// 		for {
// 			new_ch <- <-ch2
// 		}
// 	}()
// 	return new_ch
// }

// V1に比べてgoroutineが一つ少なくなる
func fanInV2(ch1, ch2 <- chan string) <- chan string {
	new_ch := make(chan string)
	go func() {
		for {
			select {
			case s := <-ch1:
				new_ch <- s
			case s := <-ch2:
				new_ch <- s
			}
		}
	}()
	return new_ch
}

// func generator(msg string) <-chan string {
// 	ch := make(chan string)
// 	go func() {
// 		for i := 0; ; i++ {
// 			ch <- fmt.Sprintf("%s, %d", msg, i)
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	return ch
// }

func generator(msg string, quit chan bool) <- chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case ch <- fmt.Sprintf("%s", msg):
				// 
			case <- quit:
				fmt.Println("Goroutine done")
				return
			}
		}
	}()
	return ch
}

func main() {
	// ch := generator("Hello")
	// for i := 0; i < 5; i++ {
	// 	fmt.Println(<-ch)
	// }

	// ch := fanInV2(generator("Hello"), generator("Bye"))
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<- ch)
	// }

	// ch := generator("Hi!")
	// for i := 0; i < 10; i++ {
	// 	select {
	// 	case s := <-ch:
	// 		fmt.Println(s)
	// 	case <-time.After(2 * time.Second): // 各ループごとにtime.Afterが初期化されて実行される
	// 		fmt.Println("Waited too long!")
	// 		return
	// 	}
	// }

	quit := make(chan bool)
	ch := generator("Hi!", quit)
	for i := rand.Intn(50); i >= 0; i-- {
		fmt.Println(<-ch, i)
	}
	quit <- true
}