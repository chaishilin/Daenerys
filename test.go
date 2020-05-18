package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	id := 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i:=0;i<10;i++{
		mu.Lock()
		fmt.Printf("pid:无, id:%d\n",id)
		id++
		time.Sleep(10*time.Microsecond)
		mu.Unlock()
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				mu.Lock()
				fmt.Printf("pid:%d, id:%d\n",i,id)
				id++
				time.Sleep(10*time.Microsecond)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("should be :",30*200,id)
}
