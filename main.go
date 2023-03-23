package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gavincabbage/kurz/api"
	"github.com/gavincabbage/kurz/store/inmem"
	kurzredis "github.com/gavincabbage/kurz/store/redis"
)

var argv struct {
	listenAddr string
	redisAddr  string
}

func init() {
	flag.StringVar(&argv.listenAddr, "listen-addr", ":8080", "host:port on which to serve the http api")
	flag.StringVar(&argv.redisAddr, "redis-addr", "", "address of redis to use as link store")
}

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

	var store api.LinkStore
	if argv.redisAddr != "" {
		cli := redis.NewClient(&redis.Options{
			Addr: argv.redisAddr,
		})
		store = kurzredis.New(cli)
	} else {
		store = &inmem.Store{}
	}

	server := api.New(logger, argv.listenAddr, store)
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
