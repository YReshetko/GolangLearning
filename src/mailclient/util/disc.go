package util

import (
	"fmt"
	"log"
	"strings"
	"syscall"
	"unsafe"
)

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	_, _, err := c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path[:2]))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "successfully") {
			log.Println("Error during retrieving memory statistic:", err)
			return
		}
	}
	disk.All = uint64(lpTotalNumberOfBytes)
	disk.Free = uint64(lpTotalNumberOfFreeBytes)
	disk.Used = disk.All - disk.Free
	return
}
