package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beratdogan6/nodeprobe/internal/collector"
)

func main() {
	f, err := os.Open("/proc/loadavg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	load, err := collector.ParseLoadAvg(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", load)
}
