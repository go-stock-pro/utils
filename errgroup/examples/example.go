package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rohanraj7316/utils/errgroup"
)

func main() {
	errGrp := errgroup.New()
	errGrp.SetLimit(10)

	timeout := 200 * time.Second

	result := []string{}

	errGrp.GoWithTimeout(context.Background(), timeout, func() error {
		for i := 0; i < 1000000000000000000; i++ {
			result = append(result, time.Since(time.Now()).String())
		}

		return nil
	})

	if err := errGrp.Wait(); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Println("end: ", len(result))
}
