package extension

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"log"

	"github.com/moloch--/memmod"
)

const (
	Success = 0
	Failure = 1
)

type WindowsExtension struct {
	id     string
	data   []byte
	module *memmod.Module
	arch   string
	init   string
	sync.Mutex
}

// NewWindowsExtension - Load a new windows extension
func NewWindowsExtension(data []byte, id string, arch string, init string) *WindowsExtension {
	return &WindowsExtension{
		id:   id,
		data: data,
		arch: arch,
		init: init,
	}
}

// GetID - Get the extension ID
func (w *WindowsExtension) GetID() string {
	return w.id
}

// GetArch - Get the extension architecture
func (w *WindowsExtension) GetArch() string {
	return w.arch
}

// Load - Load the extension into memory
func (w *WindowsExtension) Load() error {
	var err error
	if len(w.data) == 0 {
		return errors.New("extension data is empty")
	}
	w.Lock()
	defer w.Unlock()
	w.module, err = memmod.LoadLibrary(w.data)
	if err != nil {
		fmt.Print("ERROR loading")
		return err
	}
	if w.init != "" {
		initProc, errInit := w.module.ProcAddressByName(w.init)
		if errInit == nil {
			log.Printf("Calling %s\n", w.init)
			syscall.Syscall(initProc, 0, 0, 0, 0)
		} else {
			return errInit
		}
	}
	return nil
}

// Call - Call an extension export
func (w *WindowsExtension) Call(export string, arguments []byte, onFinish func([]byte)) error {
	var (
		argumentsPtr  uintptr
		argumentsSize uintptr
	)
	if w.module == nil {
		return errors.New("module not loaded")
	}
	callback := syscall.NewCallback(newWindowsExtensionCallback(onFinish))
	exportPtr, err := w.module.ProcAddressByName(export)
	if err != nil {
		fmt.Println("Export pointer: ", err)
		return err
	}
	if len(arguments) > 0 {
		argumentsPtr = uintptr(unsafe.Pointer(&arguments[0]))
		argumentsSize = uintptr(uint32(len(arguments)))
	}
	fmt.Printf("Calling %s, arguments addr: 0x%08x, args size: %08x\n", export, argumentsPtr, argumentsSize)
	// The extension API must respect the following prototype:
	// int Run(buffer char*, bufferSize uint32_t, goCallback callback)
	// where goCallback = int(char *, int)
	w.Lock()
	defer w.Unlock()
	_, _, errNo := syscall.Syscall(exportPtr, 3, argumentsPtr, argumentsSize, callback)
	if errNo != 0 {
		fmt.Println("Syscall error n: ", errNo)
		return errors.New(errNo.Error())
	}

	return nil
}

func newWindowsExtensionCallback(onFinish func([]byte)) func(data uintptr, dataLen uintptr) uintptr {
	// extensionCallback takes a buffer (char *) and its size (int) as parameters
	// so we can pass data back to the Go process from the loaded DLL
	return func(data uintptr, dataLen uintptr) uintptr {
		outDataSize := int(dataLen)
		outBytes := unsafe.Slice((*byte)(unsafe.Pointer(data)), outDataSize)
		if dataLen > 0 {
			onFinish(outBytes)
		}
		return Success
	}
}
