package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	addr := os.Getenv("addr")
	if addr == "" {
		addr = ":8000"
	}

	// start http server
	log.Printf("Starting server on: [%s]\n", addr)
	q, err := serve(addr)
	if err != nil {
		log.Printf("unable to start server: %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, done := context.WithTimeout(ctx, time.Second*10)
	defer done()

	quit(ctx, q)
}

func quit(ctx context.Context, fs ...func(context.Context) error) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	wg := sync.WaitGroup{}

	for _, f := range fs {
		wg.Add(1)

		go func(f func(ctx context.Context) error) {

			err := f(ctx)
			if err != nil {
				log.Printf("did not shutdown cleanly: %s", err)
			}

			wg.Done()
		}(f)
	}

	c := make(chan struct{})

	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		log.Println("clean shutdown")
	case <-ctx.Done():
		log.Println("deadline exceeded, forcing quit")
	}

}
