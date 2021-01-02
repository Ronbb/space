package usn

import (
	"errors"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32         = windows.NewLazySystemDLL("kernel32.dll")
	procOpenFileByID = kernel32.NewProc("OpenFileById")
)

// Volume .
func Volume(path string) string {
	for {
		last := filepath.Dir(strings.TrimRight(path, "\\"))
		if strings.HasSuffix(last, ".") {
			return path
		}
		path = last
	}
}

// IsNTFS .
func IsNTFS(volume string) (isNTFS bool, err error) {
	var maximumComponentLength uint32
	fileSystemNameBuffer := [1000]uint16{}
	err = windows.GetVolumeInformation(syscall.StringToUTF16Ptr(volume), nil, 0, nil, &maximumComponentLength, nil, &fileSystemNameBuffer[0], 1000)
	if err != nil {
		return false, err
	}

	return string(syscall.UTF16ToString(fileSystemNameBuffer[:])) == "NTFS", nil
}

// NewHandle .
func NewHandle(volume string) (syscall.Handle, error) {
	volumeName := [1000]uint16{}
	err := windows.GetVolumeNameForVolumeMountPoint(syscall.StringToUTF16Ptr(volume), &volumeName[0], 1000)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	file := syscall.UTF16ToString(volumeName[:])
	file = strings.TrimRight(file, "\\")

	h, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(file),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		0, 0,
	)
	if err != nil {
		println(err.Error())
		return 0, nil
	}
	return h, nil
}

type id struct {
	size uint32
	t uint32
	b16 [16]byte
	
}

// GetSizeByID .
func GetSizeByID(handle syscall.Handle, b16 []byte) (int64, error) {
	if len(b16) != 16 {
		return 0, errors.New("b16's length is not 16")
	}
	b := [16]byte{}
	for i := 0; i < 16; i++ {
		b[i] = b16[i]
	}
	h, _, e := procOpenFileByID.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&id{
			size: 24,
			t: 0,
			b16: b,
		})),
		uintptr(syscall.SYNCHRONIZE|0x0080),
		uintptr(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE),
		uintptr(0),
		uintptr(0),
	)
	file := syscall.Handle(h)
	if file == syscall.InvalidHandle {
		return 0, e
	}
	info := windows.ByHandleFileInformation{}
	err := windows.GetFileInformationByHandle(windows.Handle(h), &info)
	if err != nil {
		return 0, err
	}
	return int64(uint64(info.FileSizeHigh)<<32 | uint64(info.FileSizeLow)), nil
}
