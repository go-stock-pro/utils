package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pkg/errors"
	egp "golang.org/x/sync/errgroup"
)

type Handler struct {
	eg *egp.Group
}

func New() *Handler {
	eg := &egp.Group{}

	return &Handler{
		eg: eg,
	}
}

func (h *Handler) SetLimit(n int) {
	h.eg.SetLimit(n)
}

func (h *Handler) Go(fn func() error) {
	h.eg.Go(fn)
}

func (h *Handler) Wait() error {
	return h.eg.Wait()
}

func (h *Handler) GoWithTimeout(ctx context.Context, timeout time.Duration, fn func() error) {
	h.eg.Go(func() error {
		nctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		errChan := make(chan error, 1)
		panicChan := make(chan string, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			errChan <- fn()
		}()

		select {
		case pan := <-panicChan:
			return errors.New(pan)
		case err := <-errChan:
			return err
		case <-nctx.Done():
			return nctx.Err()
		}
	})
}

func (h *Handler) TryGo(fn func() error) bool {
	return h.eg.TryGo(fn)
}
