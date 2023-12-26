package computer

import (
	"C"
	"fmt"
	"syscall"
)

type (
	ATOM      uint16
	HANDLE    uintptr
	HGLOBAL   HANDLE
	HINSTANCE HANDLE
	LCID      uint32
)

// Windows API functions
var (
	modKernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetThreadLocale = modKernel32.NewProc("GetThreadLocale")
)

// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid/70feba9f-294e-491e-b6eb-56532684c37f
const (
	es_ES = 3082 //0x0C0A
	en_US = 1033
)

func threadLocale() LCID {
	ret, _, _ := procGetThreadLocale.Call()
	return LCID(ret)
}

func GetLocale() string {
	lcid := threadLocale()

	switch lcid {
	case es_ES:
		return "es-ES"
	case en_US:
		return "en-US"
	default:
		return fmt.Sprintf("Unkown: %d", lcid)
	}
}
