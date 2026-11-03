// Package metrics writes metrics in the Prometheus text exposition format.
package metrics

import (
	"fmt"
	"io"
	"strconv"
)

// ContentType is the value to send in the Content-Type header of a
// response body written by this package.
const ContentType = "text/plain; version=0.0.4; charset=utf-8"

// WriteGauge writes one gauge metric to w: a HELP line, a TYPE line and
// the sample itself. The value is formatted with the shortest decimal
// representation that round-trips.
func WriteGauge(w io.Writer, name, help string, value float64) {
	_, _ = fmt.Fprintf(w, "# HELP %s %s\n", name, help)
	_, _ = fmt.Fprintf(w, "# TYPE %s gauge\n", name)
	_, _ = fmt.Fprintf(w, "%s %s\n", name, strconv.FormatFloat(value, 'f', -1, 64))
}
