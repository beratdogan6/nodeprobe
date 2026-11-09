// Command nodeprobe serves Linux host metrics over HTTP in the
// Prometheus text format.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/beratdogan6/nodeprobe/internal/server"
)

// shutdownTimeout is how long in-flight requests get to finish after a
// signal arrives, before the server is closed out from under them.
const shutdownTimeout = 5 * time.Second

func main() {
	addr := flag.String("addr", ":9200", "address to listen on")
	flag.Parse()

	if err := run(*addr); err != nil {
		log.Fatal(err)
	}
	log.Println("shutting down")
}

// run listens on addr until SIGINT or SIGTERM arrives, then drains the
// server. It returns nil when the shutdown was clean.
func run(addr string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              addr,
		Handler:           server.New(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		log.Println("nodeprobe listening on " + addr)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		// ListenAndServe only returns on failure here; a shutdown
		// started by a signal is handled in the other case.
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		stop() // restore the default handlers, so a second signal kills us
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
