package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ServeWithGracefulShutdown(
	ctx context.Context,
	server *http.Server,
	listener net.Listener,
	shutdownTimeout time.Duration,
) error {
	// Flow: serve HTTP in a background goroutine -> wait for a serve error or context cancellation -> shut down with a timeout -> return the final error
	if server == nil {
		return errors.New("server is required")
	}
	if listener == nil {
		return errors.New("listener is required")
	}
	if shutdownTimeout <= 0 {
		return errors.New("shutdown timeout must be greater than 0")
	}

	serveErrCh := make(chan error, 1)
	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrCh <- fmt.Errorf("serve http server: %w", err)
			return
		}
		serveErrCh <- nil
	}()

	select {
	case err := <-serveErrCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		shutdownErr := server.Shutdown(shutdownCtx)
		serveErr := <-serveErrCh

		if shutdownErr != nil {
			return fmt.Errorf("shutdown http server: %w", shutdownErr)
		}
		if serveErr != nil {
			return serveErr
		}
		return nil
	}
}

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})

	return &http.Server{Handler: mux}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := ServeWithGracefulShutdown(ctx, newServer(), listener, 5*time.Second); err != nil {
		log.Fatal(err)
	}
}
