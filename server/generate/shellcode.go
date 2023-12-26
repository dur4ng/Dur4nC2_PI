package generate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Binject/go-donut/donut"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// DonutShellcodeFromPE returns a Donut shellcode for the given PE file
func DonutShellcodeFromPEBytes(pe []byte, arch string, params string, className string, method string, isDLL bool, isUnicode bool) (data []byte, err error) {
	ext := ".exe"
	if isDLL {
		ext = ".dll"
	}
	var isUnicodeVar uint32
	if isUnicode {
		isUnicodeVar = 1
	}
	//donutArch := getDonutArch(arch)

	// We don't use DonutConfig.Thread = 1 because we create our own remote thread
	// in the task runner, and we're doing some housekeeping on it.
	// Having DonutConfig.Thread = 1 means another thread will be created
	// inside the one we created, and that will fuck up our monitoring
	// since we can't grab a handle to the thread created by the Donut loader,
	// and thus the waitForCompletion call will most of the time never complete.
	config := donut.DonutConfig{
		Type:       getDonutType(ext, false),
		InstType:   donut.DONUT_INSTANCE_PIC,
		Parameters: params,
		Class:      className,
		Method:     method,
		Bypass:     1,         // 1=skip, 2=abort on fail, 3=continue on fail.
		Format:     uint32(1), // 1=raw, 2=base64, 3=c, 4=ruby, 5=python, 6=powershell, 7=C#, 8=hex
		Arch:       getDonutArch(arch),
		Entropy:    1,         // 1=disable, 2=use random names, 3=random names + symmetric encryption (default)
		Compress:   uint32(1), // 1=disable, 2=LZNT1, 3=Xpress, 4=Xpress Huffman
		ExitOpt:    1,         // exit thread
		Unicode:    isUnicodeVar,
	}
	return getDonut(pe, &config)
}

func DonutFromAssembly(assembly []byte, isDLL bool, arch string, params string, method string, runtime string, className string, appDomain string) ([]byte, error) {
	ext := ".exe"
	if isDLL {
		ext = ".dll"
	}

	config := donut.DonutConfig{
		Type:       getDonutType(ext, true),
		InstType:   donut.DONUT_INSTANCE_PIC,
		Parameters: params,
		Class:      className,
		Runtime:    runtime,
		Domain:     appDomain,
		Method:     method,
		Bypass:     3,         // 1=skip, 2=abort on fail, 3=continue on fail.
		Format:     uint32(1), // 1=raw, 2=base64, 3=c, 4=ruby, 5=python, 6=powershell, 7=C#, 8=hex
		Arch:       getDonutArch("x84"),
		Entropy:    1,         // 1=disable, 2=use random names, 3=random names + symmetric encryption (default)
		Compress:   uint32(1), // 1=disable, 2=LZNT1, 3=Xpress, 4=Xpress Huffman
		ExitOpt:    1,         // exit thread
		Unicode:    0,
	}

	return getDonut(assembly, &config)
}

func DonutShellcodeFromPath(pePath string) (data []byte, err error) {
	/*
		ext := strings.Split(pePath, ".")
		extStr := ".exe"
		if ext[1] == "dll" {
			extStr = ".dll"
		}
	*/

	config := donut.DonutConfig{
		//Type:       getDonutType(extStr, true),
		InstType:   donut.DONUT_INSTANCE_PIC,
		Parameters: "",
		Class:      "",
		Runtime:    "",
		Domain:     "",
		Method:     "",
		Bypass:     3,         // 1=skip, 2=abort on fail, 3=continue on fail.
		Format:     uint32(1), // 1=raw, 2=base64, 3=c, 4=ruby, 5=python, 6=powershell, 7=C#, 8=hex
		Arch:       getDonutArch("x84"),
		Entropy:    1,         // 1=disable, 2=use random names, 3=random names + symmetric encryption (default)
		Compress:   uint32(1), // 1=disable, 2=LZNT1, 3=Xpress, 4=Xpress Huffman
		ExitOpt:    1,         // exit thread
		Unicode:    0,
	}
	shellcode, err := donut.ShellcodeFromFile(pePath, &config)
	if err != nil {
		return nil, err
	}
	switch runtime.GOOS {
	case "windows":
		return nil, errors.New("can not generate donuts from windows")
	case "linux":
		if err == nil {
			f, err := os.Create("/tmp/loader.bin")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			if _, err = shellcode.WriteTo(f); err != nil {
				log.Fatal(err)
			}
		}
		fileBytes, err := ioutil.ReadFile("/tmp/loader.bin")
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil, nil
		}
		return fileBytes, nil
	default:
		return nil, errors.New("unsupported os")
	}
}

func CDonutShellcodeFromPath(filePath string) ([]byte, error) {
	var (
		shell string
		cmd   *exec.Cmd
	)
	switch runtime.GOOS {
	case "windows":
		return nil, errors.New("This feature is not suppported in this OS")

	case "linux":
		shell = "bash"
		cmd = exec.Command(shell)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer stdin.Close()
			execStr := `/home/dur4n/repos/donut/donut -i ` + filePath + ` -o /tmp/cdonutLoader.bin`
			fmt.Println(execStr)
			fmt.Fprintln(stdin, execStr)
			out, _ := cmd.CombinedOutput()
			fmt.Println(string(out))
		}()

		shellcode, err := ioutil.ReadFile("/tmp/cdonutLoader.bin")
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil, err
		}
		return shellcode, nil
	default:
		return nil, errors.New("target os not available")
	}
}

func getDonut(data []byte, config *donut.DonutConfig) (shellcode []byte, err error) {
	buf := bytes.NewBuffer(data)
	res, err := donut.ShellcodeFromBytes(buf, config)
	if err != nil {
		return
	}
	shellcode = res.Bytes()
	stackCheckPrologue := []byte{
		// Check stack is 8 byte but not 16 byte aligned or else errors in LoadLibrary
		0x48, 0x83, 0xE4, 0xF0, // and rsp,0xfffffffffffffff0
		0x48, 0x83, 0xC4, 0x08, // add rsp,0x8
	}
	shellcode = append(stackCheckPrologue, shellcode...)
	return
}
func getDonutArch(arch string) donut.DonutArch {
	var donutArch donut.DonutArch
	switch strings.ToLower(arch) {
	case "x32", "386":
		donutArch = donut.X32
	case "x64", "amd64":
		donutArch = donut.X64
	case "x84":
		donutArch = donut.X84
	default:
		donutArch = donut.X84
	}
	return donutArch
}
func getDonutType(ext string, dotnet bool) donut.ModuleType {
	var donutType donut.ModuleType
	switch strings.ToLower(filepath.Ext(ext)) {
	case ".exe", ".bin":
		if dotnet {
			donutType = donut.DONUT_MODULE_NET_EXE
		} else {
			donutType = donut.DONUT_MODULE_EXE
		}
	case ".dll":
		if dotnet {
			donutType = donut.DONUT_MODULE_NET_DLL
		} else {
			donutType = donut.DONUT_MODULE_DLL
		}
	case ".xsl":
		donutType = donut.DONUT_MODULE_XSL
	case ".js":
		donutType = donut.DONUT_MODULE_JS
	case ".vbs":
		donutType = donut.DONUT_MODULE_VBS
	}
	return donutType
}
