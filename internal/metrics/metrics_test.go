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

func TestWriteCounter(t *testing.T) {
	tests := []struct {
		name    string
		counter string
		help    string
		samples []Sample
		want    string
	}{
		{
			name:    "one sample line per device",
			counter: "nodeprobe_network_receive_bytes_total",
			help:    "Bytes received since boot.",
			samples: []Sample{
				{Labels: []Label{{Name: "device", Value: "lo"}}, Value: 72467},
				{Labels: []Label{{Name: "device", Value: "eth0"}}, Value: 303842249},
			},
			want: "# HELP nodeprobe_network_receive_bytes_total Bytes received since boot.\n" +
				"# TYPE nodeprobe_network_receive_bytes_total counter\n" +
				"nodeprobe_network_receive_bytes_total{device=\"lo\"} 72467\n" +
				"nodeprobe_network_receive_bytes_total{device=\"eth0\"} 303842249\n",
		},
		{
			name:    "escapes backslash, quote and newline in label values",
			counter: "nodeprobe_network_receive_bytes_total",
			help:    "Bytes received since boot.",
			samples: []Sample{
				{Labels: []Label{{Name: "device", Value: "we\\ird\"one\nhere"}}, Value: 1},
			},
			want: "# HELP nodeprobe_network_receive_bytes_total Bytes received since boot.\n" +
				"# TYPE nodeprobe_network_receive_bytes_total counter\n" +
				`nodeprobe_network_receive_bytes_total{device="we\\ird\"one\nhere"} 1` + "\n",
		},
		{
			name:    "several labels are comma separated",
			counter: "nodeprobe_network_transmit_bytes_total",
			help:    "Bytes transmitted since boot.",
			samples: []Sample{
				{Labels: []Label{
					{Name: "device", Value: "eth0"},
					{Name: "kind", Value: "physical"},
				}, Value: 3637840},
			},
			want: "# HELP nodeprobe_network_transmit_bytes_total Bytes transmitted since boot.\n" +
				"# TYPE nodeprobe_network_transmit_bytes_total counter\n" +
				"nodeprobe_network_transmit_bytes_total{device=\"eth0\",kind=\"physical\"} 3637840\n",
		},
		{
			name:    "no labels and no samples",
			counter: "nodeprobe_network_receive_bytes_total",
			help:    "Bytes received since boot.",
			samples: nil,
			want: "# HELP nodeprobe_network_receive_bytes_total Bytes received since boot.\n" +
				"# TYPE nodeprobe_network_receive_bytes_total counter\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			WriteCounter(&sb, tt.counter, tt.help, tt.samples)
			if got := sb.String(); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
