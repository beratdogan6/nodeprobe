package metrics

import (
	"strings"
	"testing"
)

func TestWriteGauge(t *testing.T) {
	tests := []struct {
		name  string
		gauge string
		help  string
		value float64
		want  string
	}{
		{
			name:  "fractional value",
			gauge: "nodeprobe_load1",
			help:  "1-minute load average.",
			value: 0.03,
			want: "# HELP nodeprobe_load1 1-minute load average.\n" +
				"# TYPE nodeprobe_load1 gauge\n" +
				"nodeprobe_load1 0.03\n",
		},
		{
			name:  "whole number keeps no decimal point",
			gauge: "nodeprobe_memory_total_bytes",
			help:  "Total usable memory.",
			value: 16322797568,
			want: "# HELP nodeprobe_memory_total_bytes Total usable memory.\n" +
				"# TYPE nodeprobe_memory_total_bytes gauge\n" +
				"nodeprobe_memory_total_bytes 16322797568\n",
		},
		{
			name:  "zero",
			gauge: "nodeprobe_disk_used_bytes",
			help:  "Used space on the root filesystem.",
			value: 0,
			want: "# HELP nodeprobe_disk_used_bytes Used space on the root filesystem.\n" +
				"# TYPE nodeprobe_disk_used_bytes gauge\n" +
				"nodeprobe_disk_used_bytes 0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			WriteGauge(&sb, tt.gauge, tt.help, tt.value)
			if got := sb.String(); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
