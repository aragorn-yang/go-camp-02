package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	appHandler := http.NewServeMux()
	appHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("go camp"))
	})

	app := http.Server{
		Handler: appHandler,
		Addr:    ":8080",
	}

	g.Go(func() error {
		return app.ListenAndServe()
	})

	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			fmt.Printf("get signal %s, application is gonna shut down", sig)
			app.Shutdown(ctx)
			return fmt.Errorf("received os signal: %v", sig)
		}
	})

	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Print("context was canceled")
		} else {
			fmt.Printf("received error: %v", err)
		}
	} else {
		fmt.Println("finished clean")
	}
}
