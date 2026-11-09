// Package metrics writes metrics in the Prometheus text exposition format.
package metrics

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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

// Label is one name/value pair attached to a sample, such as
// device="eth0". Label values are escaped when written.
type Label struct {
	Name  string
	Value string
}

// Sample is a single measurement of a metric, identified by its labels.
// A metric with no labels has exactly one sample with no labels.
type Sample struct {
	Labels []Label
	Value  float64
}

// WriteCounter writes one counter metric to w: a HELP line, a TYPE line
// and then one sample line per entry in samples, in the given order.
// The HELP and TYPE lines are written once for the metric name, which is
// what the text format requires even when several samples share it.
func WriteCounter(w io.Writer, name, help string, samples []Sample) {
	_, _ = fmt.Fprintf(w, "# HELP %s %s\n", name, help)
	_, _ = fmt.Fprintf(w, "# TYPE %s counter\n", name)
	for _, s := range samples {
		_, _ = fmt.Fprintf(w, "%s%s %s\n", name, formatLabels(s.Labels),
			strconv.FormatFloat(s.Value, 'f', -1, 64))
	}
}

// formatLabels renders labels as {name="value",...}, or as the empty
// string when there are none.
func formatLabels(labels []Label) string {
	if len(labels) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("{")
	for i, l := range labels {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(l.Name)
		sb.WriteString(`="`)
		sb.WriteString(escapeLabelValue(l.Value))
		sb.WriteString(`"`)
	}
	sb.WriteString("}")
	return sb.String()
}

// escapeLabelValue escapes the three characters the text format does not
// allow to appear literally inside a quoted label value: a backslash, a
// double quote and a line feed.
func escapeLabelValue(value string) string {
	return labelValueEscaper.Replace(value)
}

// labelValueEscaper is shared by every label value written, so the
// replacement table is built once.
var labelValueEscaper = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	"\n", `\n`,
)
