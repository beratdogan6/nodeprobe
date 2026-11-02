package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beratdogan6/nodeprobe/internal/collector"
)

func main() {
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
}
