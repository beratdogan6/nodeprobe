// Command nodeprobe serves Linux host metrics over HTTP in the
// Prometheus text format.
package main

import (
	"log"
	"net/http"

	"github.com/beratdogan6/nodeprobe/internal/server"
)

// addr is the address nodeprobe listens on.
const addr = ":9200"

func main() {
	log.Println("nodeprobe listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, server.New()))
}
