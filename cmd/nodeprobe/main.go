package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/beratdogan6/nodeprobe/internal/collector"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "ok")
	})

	log.Fatal(http.ListenAndServe(":9200", mux))

	loadFile, err := os.Open("/proc/loadavg")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = loadFile.Close() }()

	load, err := collector.ParseLoadAvg(loadFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", load)

	memFile, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = memFile.Close() }()

	mem, err := collector.ParseMemInfo(memFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", mem)

	disk, err := collector.ReadDisk("/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", disk)
}
