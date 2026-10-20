package collector

import (
	"strings"
	"testing"
)

func TestParseLoadAvg(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    LoadAvg
		wantErr bool
	}{
		{
			name:  "valid",
			input: "0.03 0.11 0.08 1/238 6718\n",
			want:  LoadAvg{Load1: 0.03, Load5: 0.11, Load15: 0.08},
		},
		{
			name:    "too few fields",
			input:   "0.03 0.11",
			wantErr: true,
		},
		{
			name:    "not a number",
			input:   "abc 0.11 0.08 1/238 6718",
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
			got, err := ParseLoadAvg(strings.NewReader(tt.input))

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
