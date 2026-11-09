package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// netDevFields is the number of counters each /proc/net/dev line carries
// after the interface name: 8 receive columns followed by 8 transmit ones.
const netDevFields = 16

// NetDev holds the byte counters of one network interface, as read from
// /proc/net/dev. Both counters are totals since boot, not rates.
type NetDev struct {
	Name          string
	ReceiveBytes  uint64
	TransmitBytes uint64
}

// ReadNetDev opens /proc/net/dev and parses it.
func ReadNetDev() ([]NetDev, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return ParseNetDev(f)
}

// ParseNetDev reads /proc/net/dev-formatted data from r and returns one
// NetDev per interface, in the order the file lists them. Lines that do
// not look like interface rows are skipped. It returns an error if
// reading fails or if no interface could be parsed at all.
func ParseNetDev(r io.Reader) ([]NetDev, error) {
	var devices []NetDev

	scanner := bufio.NewScanner(r)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		if lineNo < 2 { // the two header lines
			continue
		}

		// The name is glued to the colon, but the first counter may be
		// glued to it too: "eth0:3038422490000 ...".
		name, rest, ok := strings.Cut(scanner.Text(), ":")
		if !ok {
			continue
		}

		name = strings.TrimSpace(name)
		fields := strings.Fields(rest)
		if name == "" || len(fields) < netDevFields {
			continue
		}

		receive, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}
		transmit, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			continue
		}

		devices = append(devices, NetDev{
			Name:          name,
			ReceiveBytes:  receive,
			TransmitBytes: transmit,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan net/dev: %w", err)
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("no interfaces found")
	}

	return devices, nil
}
