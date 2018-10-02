package util

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func OpenWindow(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start() == nil
}

func IsRelativePath(path string) bool {
	return string(path[1]) != ":"
}

func CreateAbsolutePath(relativePath string) string {
	pathToRoot, err := filepath.Abs(".")
	if err != nil {
		log.Println("Error during retrieving absolute root path of application:", err)
		return relativePath
	} else {
		if string(relativePath[0]) != "/" {
			return pathToRoot + "/" + relativePath
		} else {
			return pathToRoot + relativePath
		}

	}
}
func RunWinProgramm(program string, args []string) bool {
	cmd := exec.Command(program, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start() == nil
}

func KillPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	return nil
}

// TH32CS_SNAPPROCESS is described in https://msdn.microsoft.com/de-de/library/windows/desktop/ms682489(v=vs.85).aspx
const TH32CS_SNAPPROCESS = 0x00000002

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func FindPIDByName(name string) (int, error) {
	procs, err := processes()
	if err != nil {
		log.Println("Error Win processes search:", err)
		return 0, err
	}
	explorer := findProcessByName(procs, name)
	if explorer != nil {
		// found it
		return explorer.ProcessID, nil
	}
	return 0, nil
}

func processes() ([]WindowsProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}
