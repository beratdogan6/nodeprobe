package collector

import (
	"fmt"
	"syscall"
)

// DiskInfo holds filesystem usage for a mount point, in bytes.
type DiskInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Used      uint64
}

// ReadDisk asks the kernel (statfs syscall) for usage of the filesystem
// that contains path, e.g. "/".
func ReadDisk(path string) (DiskInfo, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return DiskInfo{}, fmt.Errorf("statfs %s: %w", path, err)
	}
	return diskFromStatfs(st), nil
}

// diskFromStatfs converts raw statfs block counts into bytes.
func diskFromStatfs(st syscall.Statfs_t) DiskInfo {
	bsize := uint64(st.Bsize)
	return DiskInfo{
		Total:     st.Blocks * bsize,
		Free:      st.Bfree * bsize,
		Available: st.Bavail * bsize,
		Used:      (st.Blocks - st.Bfree) * bsize,
	}
}
