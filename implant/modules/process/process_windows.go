package process

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}

// Windows API functions
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	modadvapi32                  = syscall.NewLazyDLL("advapi32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procOpenProcess              = modKernel32.NewProc("OpenProcess")
	procOpenProcessToken         = modadvapi32.NewProc("OpenProcessToken")
	procGetTokenInformation      = modadvapi32.NewProc("GetTokenInformation")
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
)

// PROCESSENTRY32 is the Windows API structure that contains a process's
// information.
type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	pid  int
	ppid int
	exe  string
}

func (p *WindowsProcess) Pid() int {
	return p.pid
}

func (p *WindowsProcess) PPid() int {
	return p.ppid
}

func (p *WindowsProcess) Executable() string {
	return p.exe
}

func newWindowsProcess(e *PROCESSENTRY32) *WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return &WindowsProcess{
		pid:  int(e.ProcessID),
		ppid: int(e.ParentProcessID),
		exe:  syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func FindProcess(pid int) (Process, error) {
	ps, err := processes()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.Pid() == pid {
			return p, nil
		}
	}

	return nil, nil
}

func processes() ([]Process, error) {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002,
		0)
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Error retrieving process info.")
	}

	results := make([]Process, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return results, nil
}

func IsProcessCurrentElevated() bool {
	//var accessRights uint32 = windows.READ_CONTROL
	//_syscall.OpenProcess(accessRights, 0, )
	//handle, _ := windows.GetCurrentProcess()
	token, _ := windows.OpenCurrentProcessToken()
	defer token.Close()

	//token.IsElevated()
	/*
		n := uint32(0)
		windows.GetTokenInformation(token, windows.TokenElevationType, nil, 0, &n)
		b := make([]byte, n)
		windows.GetTokenInformation(token, windows.TokenElevationType, &b[0], uint32(len(b)), &n)
		fmt.Scanln(&b)
	*/
	n := uint32(0)
	procGetTokenInformation.Call(uintptr(token), windows.TokenElevationType, uintptr(unsafe.Pointer(nil)), 0, uintptr(unsafe.Pointer(&n)))
	b := make([]byte, n)
	procGetTokenInformation.Call(uintptr(token), windows.TokenElevationType, uintptr(unsafe.Pointer(&b[0])), uintptr(uint32(len(b))), uintptr(unsafe.Pointer(&n)))

	switch int(b[0]) {
	case 2:
		return true
	case 3:
		return false
	default:
		return false
	}
}
