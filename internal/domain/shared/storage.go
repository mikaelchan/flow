package shared

import (
	"os"
	"syscall"
)

type StorageType uint8

const (
	Local StorageType = iota
	// Maybe add more types in the future
)

type StorageLocation struct {
	Path        string      `json:"path"`
	StorageType StorageType `json:"storage_type"`
	Capacity    uint64      `json:"capacity"` // 0 means unlimited, in GiB
}

func (sl StorageLocation) IsValid() bool {
	return sl.Path != "" && sl.StorageType.IsValid()
}

func (sl StorageLocation) IsWritable() (isWritable bool, err error) {
	if sl.StorageType == Local {
		isWritable = false
		info, err := os.Stat(sl.Path)
		if err != nil {
			return isWritable, err
		}
		if !info.IsDir() {
			return isWritable, nil
		}
		if info.Mode().Perm()&(1<<(uint(7))) == 0 {
			return isWritable, nil
		}
		var stat syscall.Stat_t
		if err := syscall.Stat(sl.Path, &stat); err != nil {
			return isWritable, err
		}
		if stat.Uid != uint32(os.Getuid()) {
			return isWritable, nil
		}
	}
	return isWritable, nil
}

func (sl StorageLocation) FreeSpace() (uint64, error) {
	if sl.StorageType == Local {
		info, err := os.Stat(sl.Path)
		if err != nil {
			return 0, err
		}
		return uint64(info.Size()), nil
	}
	return 0, nil
}

func (st StorageType) IsValid() bool {
	return st == Local
}

func (st StorageType) String() string {
	return [...]string{"local"}[st]
}
