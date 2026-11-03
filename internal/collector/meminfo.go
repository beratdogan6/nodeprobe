package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// MemInfo holds memory stats from /proc/meminfo, in bytes.
type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
}

// ReadMemInfo opens /proc/meminfo and parses it.
func ReadMemInfo() (MemInfo, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemInfo{}, err
	}
	defer func() { _ = f.Close() }()
	return ParseMemInfo(f)
}

// ParseMemInfo reads /proc/meminfo-formatted data from r and returns
// total, free and available memory in bytes. It returns an error if
// any of these fields is missing.
func ParseMemInfo(r io.Reader) (MemInfo, error) {
	values := map[string]uint64{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")        // "MemTotal:" -> "MemTotal"
		val, err := strconv.ParseUint(fields[1], 10, 64) // "15940232" -> 15940232
		if err != nil {
			continue
		}

		values[key] = val
	}

	if err := scanner.Err(); err != nil {
		return MemInfo{}, fmt.Errorf("scan meminfo: %w", err)
	}

	total, ok := values["MemTotal"]
	if !ok {
		return MemInfo{}, fmt.Errorf("MemTotal not found")
	}
	free, ok := values["MemFree"]
	if !ok {
		return MemInfo{}, fmt.Errorf("MemFree not found")
	}
	available, ok := values["MemAvailable"]
	if !ok {
		return MemInfo{}, fmt.Errorf("MemAvailable not found")
	}

	return MemInfo{
		Total:     total * 1024,
		Free:      free * 1024,
		Available: available * 1024,
	}, nil
}
