package collector

import (
	"strings"
	"testing"
)

func TestParseMemInfo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    MemInfo
		wantErr bool
	}{
		{
			name: "valid",
			input: `MemTotal:       1000 kB
					MemFree:         400 kB
					MemAvailable:    700 kB
					Cached:          200 kB`,
			want: MemInfo{
				Total:     1000 * 1024,
				Free:      400 * 1024,
				Available: 700 * 1024,
			},
		},
		{
			name: "missing MemAvailable",
			input: `MemTotal:       1000 kB
					MemFree:         400 kB`,
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name: "skips garbage lines",
			input: `MemTotal:       1000 kB
					garbage
					Foo:             abc kB
					MemFree:         400 kB
					MemAvailable:    700 kB
					Cached:          200 kB`,
			want: MemInfo{
				Total:     1000 * 1024,
				Free:      400 * 1024,
				Available: 700 * 1024,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemInfo(strings.NewReader(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
