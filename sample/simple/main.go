package main

import (
	"fmt"
	"time"

	"github.com/RangelReale/ecapplog-go"
)

func main() {
	c := ecapplog.NewClient(ecapplog.WithAppName("ECALGO-SAMPLE"))
	c.Open()
	defer c.Close()

	for i := 0; i < 30; i++ {
		c.Log(time.Now(), ecapplog.Priority_DEBUG, "app", fmt.Sprintf("First log: %d", i),
			ecapplog.WithOriginalCategory("app.internal"))
		c.Log(time.Now(), ecapplog.Priority_INFORMATION, "app", fmt.Sprintf("Second log: %d", i))
		c.Log(time.Now(), ecapplog.Priority_ERROR, "app", fmt.Sprintf("Third log: %d", i),
			ecapplog.WithExtraCategories([]string{"app_third"}))
	}

	fmt.Printf("Waiting 15 seconds so connection isn't closed...\n")
	select {
	case <-time.After(time.Second * 15):
		break
	}
}
