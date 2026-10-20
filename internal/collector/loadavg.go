// Package collector reads system metrics from the Linux /proc filesystem.
package collector

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// LoadAvg holds the 1, 5 and 15 minute load averages from /proc/loadavg.
type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// ParseLoadAvg reads /proc/loadavg-formatted data from r and returns
// the three load averages. It returns an error if the input has fewer
// than 3 fields or if any of them is not a valid number.
func ParseLoadAvg(r io.Reader) (LoadAvg, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return LoadAvg{}, fmt.Errorf("read loadavg: %w", err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return LoadAvg{}, fmt.Errorf("expected at least 3 fields, got %d", len(fields))
	}

	var loadAvg LoadAvg

	loadAvg.Load1, err = strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return LoadAvg{}, fmt.Errorf("parse load1: %w", err)
	}

	loadAvg.Load5, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return LoadAvg{}, fmt.Errorf("parse load5: %w", err)
	}

	loadAvg.Load15, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return LoadAvg{}, fmt.Errorf("parse load15: %w", err)
	}

	return loadAvg, nil
}
