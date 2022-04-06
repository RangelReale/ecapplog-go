package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/RangelReale/ecapplog-go"
)

func main() {
	c := ecapplog.NewClient(ecapplog.WithAppName("ECALGO-SAMPLE"))
	c.Open()
	defer c.Close()

	var w sync.WaitGroup
	w.Add(100 * 3)

	fmt.Printf("Sending logs...\n")

	go func() {
		for i := 0; i < 100; i++ {
			c.Log(time.Now(), ecapplog.Priority_DEBUG, "app", fmt.Sprintf("First log: %d", i),
				ecapplog.WithOriginalCategory("app.internal"))
			w.Done()
			c.Log(time.Now(), ecapplog.Priority_INFORMATION, "app", fmt.Sprintf("Second log: %d", i))
			w.Done()
			c.Log(time.Now(), ecapplog.Priority_ERROR, "app", fmt.Sprintf("Third log: %d", i),
				ecapplog.WithExtraCategories([]string{"app_third"}))
			w.Done()

			time.Sleep(time.Millisecond * 50)
		}
	}()

	w.Wait()

	fmt.Printf("Sleeping 5 seconds...\n")
	select {
	case <-time.After(time.Second * 5):
		break
	}

	fmt.Printf("Closing and sleeping 5 seconds...\n")
	c.Close()
	select {
	case <-time.After(time.Second * 5):
		break
	}
}
