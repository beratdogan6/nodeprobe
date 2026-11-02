package collector

import (
	"syscall"
	"testing"
)

func TestDiskFromStatfs(t *testing.T) {
	tests := []struct {
		name string
		st   syscall.Statfs_t
		want DiskInfo
	}{
		{
			name: "normal",
			st:   syscall.Statfs_t{Blocks: 100, Bfree: 40, Bavail: 30, Bsize: 4096},
			want: DiskInfo{
				Total:     100 * 4096,
				Free:      40 * 4096,
				Available: 30 * 4096,
				Used:      60 * 4096,
			},
		},
		{
			name: "empty disk",
			st:   syscall.Statfs_t{Blocks: 100, Bfree: 100, Bavail: 95, Bsize: 4096},
			want: DiskInfo{
				Total:     100 * 4096,
				Free:      100 * 4096,
				Available: 95 * 4096,
				Used:      0,
			},
		},
		{
			name: "full for users, root reserve left",
			st:   syscall.Statfs_t{Blocks: 100, Bfree: 5, Bavail: 0, Bsize: 4096},
			want: DiskInfo{
				Total:     100 * 4096,
				Free:      5 * 4096,
				Available: 0,
				Used:      95 * 4096,
			},
		},
		{
			name: "different block size",
			st:   syscall.Statfs_t{Blocks: 100, Bfree: 40, Bavail: 30, Bsize: 1024},
			want: DiskInfo{
				Total:     100 * 1024,
				Free:      40 * 1024,
				Available: 30 * 1024,
				Used:      60 * 1024,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diskFromStatfs(tt.st)
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReadDiskRoot(t *testing.T) {
	got, err := ReadDisk("/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Total == 0 {
		t.Errorf("expected non-zero total, got %+v", got)
	}
}

func TestReadDiskMissingPath(t *testing.T) {
	if _, err := ReadDisk("/this/path/does/not/exist"); err == nil {
		t.Fatal("expected error, got nil")
	}
}
