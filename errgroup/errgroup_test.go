package errgroup_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rohanraj7316/utils/errgroup"
)

func Test_ConcurrentExec(t *testing.T) {
	errGrp := errgroup.New()
	timeout := 5 * time.Second

	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
	}

	for _, url := range urls {
		url := url
		errGrp.GoWithTimeout(context.Background(), timeout, func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}

			return err
		})
	}

	if err := errGrp.Wait(); err != nil {
		t.Errorf("failed to fetch urls: %s", err)
	}
}

func Test_FailedConcurrentExec(t *testing.T) {
	errGrp := errgroup.New()
	timeout := 5 * time.Second

	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}

	for _, url := range urls {
		url := url
		errGrp.GoWithTimeout(context.Background(), timeout, func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}

			return err
		})
	}

	if err := errGrp.Wait(); err == nil {
		t.Errorf("failed to fetch all urls")
	}
}

func Test_Timeout(t *testing.T) {
	errGrp := errgroup.New()
	errGrp.SetLimit(10)

	timeout := 200 * time.Millisecond

	result := []string{}

	errGrp.GoWithTimeout(context.Background(), timeout, func() error {
		for i := 0; i < 1000000000000000000; i++ {
			result = append(result, time.Since(time.Now()).String())
		}

		return nil
	})

	if err := errGrp.Wait(); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("want %s have %s", context.DeadlineExceeded, err.Error())
		}
	}
}

func Test_PanicHandling(t *testing.T) {
	errGrp := errgroup.New()

	timeout := 200 * time.Millisecond

	for i := 0; i < 10; i++ {
		i := i
		errGrp.GoWithTimeout(context.Background(), timeout, func() error {
			if i == 5 || i == 6 {
				panic(fmt.Sprintf("panic for counter: %d", i))
			}

			return nil
		})
	}

	if err := errGrp.Wait(); err == nil {
		t.Errorf("failed to exec test case: %s", err)
	}
}

func BenchmarkGo(b *testing.B) {
	fn := func() {}
	g := errgroup.New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.Go(func() error { fn(); return nil })
	}
	g.Wait()
}
