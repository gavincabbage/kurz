package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gavincabbage/kurz/api"
	"github.com/gavincabbage/kurz/store/inmem"
)

func main() {
	logger := log.New(os.Stdout, "kurzd: ", log.LstdFlags)

	terminated := make(chan os.Signal, 1)
	signal.Notify(terminated, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(terminated)
		close(terminated)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-terminated
		cancel()
	}()

	addr := ":8080"
	if a := os.Getenv("PORT"); a != "" {
		i, err := strconv.Atoi(a)
		if err != nil {
			logger.Printf("invalid port %s\n", a)
		}

		addr = fmt.Sprintf(":%d", i)
	}

	// TODO(gavincabbage): Enable other types of store.
	store := inmem.Store{}
	server := api.New(logger, addr, &store)

	serverError := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverError <- err
		}
	}()

	select {
	case err := <-serverError:
		logger.Println("server error", err)
		cancel()
	case <-ctx.Done():
	}

	timeout, forceShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer forceShutdown()

	if err := server.Shutdown(timeout); err != nil {
		logger.Fatal(err)
	}
}
