package user

import (
	"os/user"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//https://github.com/iamacarpet/go-win64api/blob/master/users.go
//https://github.com/winlabs/gowin32/blob/master/security.go
//https://github.com/trustedsec/CS-Situational-Awareness-BOF/blob/master/src/SA/whoami/entry.c

const (
	UNLEN = 256
)

// types
/* Refereces
- windows c: https://docs.microsoft.com/en-us/windows/desktop/WinProg/windows-data-types
- syscall
- go.sys
- https://github.com/AllenDang/w32/blob/master/typedef.go
*/
type (
	BOOL          uint32
	BOOLEAN       byte
	BYTE          byte
	DWORD         uint32
	LPDWORD       *DWORD
	DWORD64       uint64
	HANDLE        uintptr
	HLOCAL        uintptr
	LARGE_INTEGER int64
	LONG          int32
	LPVOID        uintptr
	SIZE_T        uintptr
	UINT          uint32
	ULONG_PTR     uintptr
	ULONGLONG     uint64
	WORD          uint16
)

// strings

// For 16-bit unicode like LPWSTR or LPCWSTR, use syscall.UTF16PtrFromString("")
var UTF16Encoded, _ = windows.UTF16PtrFromString("Example")

//For 8-bit strings like LPSTR or LPCSTR, use StringToCharPtr
/*
//https://medium.com/@justen.walker/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
func StringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}
*/
var ANSIEncoded, _ = windows.ByteSliceFromString("Example")

/*
	import "unicode/utf16"

	// StringToCharPtr converts a Go string into pointer to a null-terminated cstring.
	// This assumes the go string is already ANSI encoded.
	func StringToCharPtr(str string) *uint8 {
		chars := append([]byte(str), 0) // null terminated
		return &chars[0]
	}

	// StringToUTF16Ptr converts a Go string into a pointer to a null-terminated UTF-16 wide string.
	// This assumes str is of a UTF-8 compatible encoding so that it can be re-encoded as UTF-16.
	func StringToUTF16Ptr(str string) *uint16 {
		wchars := utf16.Encode([]rune(str + "\x00"))
		return &wchars[0]
	}
*/

// Errors
/*
if err != syscall.Errno(0) {
		return 0, err
}
syscall.GetLastError()
*/

// Windows API functions
var (
	// NewLazyDLL(load when require) vs LoadLibrary(start loading) vs windows.NewLazySystemDLL(ensure that dll search path is constrained to the windows system directory)
	modAdvapi32     = syscall.NewLazyDLL("Advapi32.dll")
	modsecur32      = syscall.NewLazyDLL("Secur32.dll")
	modAdvapi32_alt = windows.NewLazyDLL("Secur32.dll")
	// calling procedures
	procGetUserNameA   = modAdvapi32.NewProc("GetUserNameA")
	procGetUserNameExW = modsecur32.NewProc("GetUserNameExW")
)

/* syscall vs windows
 */

/* calls examples

BOOL SystemParametersInfoW(
UINT  uiAction,
UINT  uiParam,
PVOID pvParam,
UINT  fWinIni
);

for UINT we can use a decimal or a hex: 0x0014 == 20
PVOID is a pointer to any type

imagePath, _ := windows.UTF16PtrFromString(`C:\Users\User\Pictures\image.jpg`)
imagePath, _ := syscall.UTF16PtrFromString(`C:\Users\User\Pictures\image.jpg`)

procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001A)
//we can use syscall.Syscall*
syscall.Syscall : less than 4 arguments
syscall.Syscall6: 4 to 6 arguments
syscall.Syscall9: 7 to 9 arguments
syscall.Syscall12: 10 to 12 arguments
syscall.Syscall15: 13 to 15 arguments
syscall.Syscall6(procSystemParamInfo.Addr(), 20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001A)
*/

// syscall arguments
/*
	All args must be casted to a a uintptr.

	Since Go is garbage collected, Standard Go pointers do not directly point to a place in memory. The go runtime is free to change the physical memory location pointed at by a Go pointer, such as when it grows the stack. When a pointer is converted into a raw uintptr via unsafe.Pointer — it becomes just a number untracked by the Go runtime. That number may or may not point to a valid location in memory like it once did, even after the very next instruction!
	Because of this, you have to call Syscalls with pointers to memory in certain way. By using the uintptr(unsafe.Pointer(&x)) construction in the argument list
*/

// function return Raw memory

//tools
// - https://pkg.go.dev/golang.org/x/sys/windows/mkwinsyscall
// -

// !!! IMPORTANT functions that return a buffer size
/*
	When we are using a function that one of the args is an out with the buffer size, we must call this function twice, one to retriave the buff size and the final one
	n := uint32(0)
	windows.GetTokenInformation(t.token, windows.TokenPrivileges, nil, 0, &n)

	b := make([]byte, n)
	if err := windows.GetTokenInformation(t.token, windows.TokenPrivileges, &b[0], uint32(len(b)), &n); err != nil {
		return nil, err
	}

*/

func GetWhoami() string {
	//https://stackoverflow.com/questions/35238933/golang-call-getvolumeinformation-winapi-function
	//https://stackoverflow.com/questions/57068772/unexpected-fault-address-when-calling-enumprocessmodules

	/*
			var buf = make([]uint16, UNLEN+1)
		var bufSize = uint32(len(buf))
		var buf_str string
		success, _, err := procGetUserNameA.Call(
			uintptr(2),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(bufSize),
		)

		success, _, err := procGetUserNameA.Call(
			uintptr(unsafe.Pointer(StringToCharPtr(buf_str))),
			uintptr(bufSize),
		)
		if err != syscall.Errno(0) {
			log.Fatalf("Error!!!!!!!!!!!!!!1")
		}
		//syscall.Syscall(procGetUserNameA.Addr(), uintptr(2), uintptr(unsafe.Pointer(StringToCharPtr(buf_str))), uintptr(bufSize))
		fmt.Println(syscall.GetLastError())
		if success == 0 {
			log.Fatalf("Error!!!!!!!!!!!!!!2")
		}
	*/
	var nameFormat uint32 = 2

	var nameBuffre string
	var nSize uint32 = UNLEN
	/*
		r1, _, _ := syscall.Syscall(procGetUserNameExW.Addr(), 3, uintptr(nameFormat), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(nameBuffre))), uintptr(unsafe.Pointer(&nSize)))
		if r1&0xff == 0 {
			//err = errnoErr(e1)
			fmt.Errorf("Errorr!!!\n")
		}
	*/
	procGetUserNameExW.Call(uintptr(nameFormat), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(nameBuffre))), uintptr(unsafe.Pointer(&nSize)))
	//fmt.Printf("Username: %s\n", nameBuffre)
	username, _ := user.Current()
	return username.Username
}

func StringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}

/*
	func getBaseAddress(handle uintptr) int64 {
		// GetProcessImageFileNameA
		var imageFileName [200]byte
		var fileSize uint32 = 200
		var fileName string

		ret, _, _ := procGetProcessImageFileNameA.Call(handle, uintptr(unsafe.Pointer(&imageFileName)), uintptr(fileSize))

		for _, char := range imageFileName {
			if char == 0 {
				break
			}

			fileName += string(char)
		}

		fileName = fileName[24:]

		// EnumProcessModules
		moduleHandles := make([]uintptr, 1024)
		var needed int32
		const handleSize = unsafe.Sizeof(moduleHandles[0])

		ret, _, _ = procEnumProcessModules.Call(uintptr(handle), uintptr(unsafe.Pointer(&moduleHandles[0])), handleSize*uintptr(len(moduleHandles)), uintptr(unsafe.Pointer(&needed)))

		// GetModuleFileNameExA
		var finalModuleHandle uintptr

		for _, moduleHandle := range moduleHandles {
			if moduleHandle > 0 {
				var moduleFileName [200]byte
				var moduleSize uint32 = 200
				var moduleName string

				ret, _, _ = procGetModuleFileNameExA.Call(handle, uintptr(moduleHandle), uintptr(unsafe.Pointer(&moduleFileName)), uintptr(moduleSize))

				if ret != 0 {
					for _, char := range moduleFileName {
						if char == 0 {
							break
						}

						moduleName += string(char)
					}

					moduleName = moduleName[3:]

					if moduleName == fileName {
						finalModuleHandle = uintptr(moduleHandle)
						break
					}
				}
			}
		}

		return int64(finalModuleHandle)
	}
*/
