// +build windows

package clr

import (
	"golang.org/x/sys/windows"
)

var (
	CLSID_CLRMetaHost    = windows.GUID{Data1: 0x9280188d, Data2: 0x0e8e, Data3: 0x4867, Data4: [8]byte{0xb3, 0x0c, 0x7f, 0xa8, 0x38, 0x84, 0xe8, 0xde}}
	IID_ICLRMetaHost     = windows.GUID{Data1: 0xD332DB9E, Data2: 0xB9B3, Data3: 0x4125, Data4: [8]byte{0x82, 0x07, 0xA1, 0x48, 0x84, 0xF5, 0x32, 0x16}}
	IID_ICLRRuntimeInfo  = windows.GUID{Data1: 0xBD39D1D2, Data2: 0xBA2F, Data3: 0x486a, Data4: [8]byte{0x89, 0xB0, 0xB4, 0xB0, 0xCB, 0x46, 0x68, 0x91}}
	CLSID_CLRRuntimeHost = windows.GUID{Data1: 0x90F1A06E, Data2: 0x7712, Data3: 0x4762, Data4: [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}
	IID_ICLRRuntimeHost  = windows.GUID{Data1: 0x90F1A06C, Data2: 0x7712, Data3: 0x4762, Data4: [8]byte{0x86, 0xB5, 0x7A, 0x5E, 0xBA, 0x6B, 0xDB, 0x02}}
	IID_ICorRuntimeHost  = windows.GUID{Data1: 0xcb2f6722, Data2: 0xab3a, Data3: 0x11d2, Data4: [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}
	CLSID_CorRuntimeHost = windows.GUID{Data1: 0xcb2f6723, Data2: 0xab3a, Data3: 0x11d2, Data4: [8]byte{0x9c, 0x40, 0x00, 0xc0, 0x4f, 0xa3, 0x0a, 0x3e}}
	IID_AppDomain        = windows.GUID{Data1: 0x05f696dc, Data2: 0x2b29, Data3: 0x3663, Data4: [8]byte{0xad, 0x8b, 0xc4, 0x38, 0x9c, 0xf2, 0xa7, 0x13}}
	// IID_IErrorInfo is the interface ID for the Error interface 1CF2B120-547D-101B-8E65-08002B2BD119
	IID_IErrorInfo = windows.GUID{Data1: 0x1cf2b120, Data2: 0x547d, Data3: 0x101b, Data4: [8]byte{0x8e, 0x65, 0x08, 0x00, 0x2b, 0x2b, 0xd1, 0x19}}
	// DF0B3D60-548F-101B-8E65-08002B2BD119 https://docs.microsoft.com/en-us/windows/win32/api/oaidl/nn-oaidl-isupporterrorinfo
	IID_ISupportErrorInfo = windows.GUID{Data1: 0xDF0B3D60, Data2: 0x548F, Data3: 0x101B, Data4: [8]byte{0x8e, 0x65, 0x08, 0x00, 0x2b, 0x2b, 0xd1, 0x19}}
)
