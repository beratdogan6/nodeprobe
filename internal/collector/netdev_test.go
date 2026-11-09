package collector

import (
	"slices"
	"strings"
	"testing"
)

const netDevHeader = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
`

func TestParseNetDev(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []NetDev
		wantErr bool
	}{
		{
			name: "two interfaces keep file order",
			input: netDevHeader +
				"    lo:   72467     244    0    0    0     0          0         0    72467     244    0    0    0     0       0          0\n" +
				"  eth0: 303842249   93649    0    0    0     0          0        48  3637840   44760    0    0    0     0       0          0\n",
			want: []NetDev{
				{Name: "lo", ReceiveBytes: 72467, TransmitBytes: 72467},
				{Name: "eth0", ReceiveBytes: 303842249, TransmitBytes: 3637840},
			},
		},
		{
			name: "name glued to the first counter",
			input: netDevHeader +
				"eth0:3038422490000   93649    0    0    0     0          0        48  3637840   44760    0    0    0     0       0          0\n",
			want: []NetDev{
				{Name: "eth0", ReceiveBytes: 3038422490000, TransmitBytes: 3637840},
			},
		},
		{
			name: "skips malformed lines",
			input: netDevHeader +
				"garbage without a colon\n" +
				"  eth1:   10   2    0    0    0     0          0         0\n" + // too few fields
				"  eth2:  abc   2    0    0    0     0          0         0    5   2    0    0    0     0       0          0\n" + // non-numeric receive
				"  eth3:   10   2    0    0    0     0          0         0  xyz   2    0    0    0     0       0          0\n" + // non-numeric transmit
				"    lo:   72467     244    0    0    0     0          0         0    72467     244    0    0    0     0       0          0\n",
			want: []NetDev{
				{Name: "lo", ReceiveBytes: 72467, TransmitBytes: 72467},
			},
		},
		{
			name:    "header only",
			input:   netDevHeader,
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNetDev(strings.NewReader(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
