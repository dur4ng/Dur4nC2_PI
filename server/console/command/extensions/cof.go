package extensions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/desertbit/grumble"
	"golang.org/x/text/encoding/unicode"
	"io/ioutil"
	"strconv"
)

type BOFArgsBuffer struct {
	Buffer *bytes.Buffer
}

const (
	CoffLoaderName = "coff-loader"
	bofEntryPoint  = "go"
)

func GetBOFArgs(ctx *grumble.Context, binPath string, ext *ExtensionManifest) ([]byte, error) {
	var extensionArgs []byte

	//TODO remove this
	//fmt.Println(binPath)
	binData, err := ioutil.ReadFile(binPath)
	if err != nil {
		fmt.Println("Could not load the binary")
		return nil, err
	}
	argumentsStr := ctx.Args.StringList("arguments")
	argsBuffer := BOFArgsBuffer{
		Buffer: new(bytes.Buffer),
	}
	// Parse BOF arguments from grumble
	for i, arg := range ext.Arguments {
		switch arg.Type {
		case "integer":
			fallthrough
		case "int":
			//val := ctx.Args.Int(arg.Name)
			val, _ := strconv.Atoi(argumentsStr[i])
			err = argsBuffer.AddInt(uint32(val))
			if err != nil {
				return nil, err
			}
		case "short":
			//val := ctx.Args.Int(arg.Name)
			val, _ := strconv.Atoi(argumentsStr[i])
			err = argsBuffer.AddShort(uint16(val))
			if err != nil {
				return nil, err
			}
		case "string":
			//val := ctx.Args.String(arg.Name)
			val := argumentsStr[i]
			err = argsBuffer.AddString(val)
			if err != nil {
				return nil, err
			}
		case "wstring":
			//val := ctx.Args.String(arg.Name)
			val := argumentsStr[i]
			err = argsBuffer.AddWString(val)
			if err != nil {
				return nil, err
			}
		// Adding support for filepaths so we can
		// send binary data like shellcodes to BOFs
		case "file":
			val := ctx.Args.String(arg.Name)
			data, err := ioutil.ReadFile(val)
			if err != nil {
				return nil, err
			}
			err = argsBuffer.AddData(data)
			if err != nil {
				return nil, err
			}
		}
	}
	parsedArgs, err := argsBuffer.GetBuffer()
	if err != nil {
		return nil, err
	}
	// Now build the extension's argument buffer
	extensionArgsBuffer := BOFArgsBuffer{
		Buffer: new(bytes.Buffer),
	}
	//Entrypoint?
	err = extensionArgsBuffer.AddString(bofEntryPoint)
	if err != nil {
		return nil, err
	}
	err = extensionArgsBuffer.AddData(binData)
	if err != nil {
		return nil, err
	}
	err = extensionArgsBuffer.AddData(parsedArgs)
	if err != nil {
		return nil, err
	}
	extensionArgs, err = extensionArgsBuffer.GetBuffer()
	if err != nil {
		return nil, err
	}
	return extensionArgs, nil
}

func (b *BOFArgsBuffer) AddData(d []byte) error {
	dataLen := uint32(len(d))
	err := binary.Write(b.Buffer, binary.LittleEndian, &dataLen)
	if err != nil {
		return err
	}
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}

func (b *BOFArgsBuffer) AddShort(d uint16) error {
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}

func (b *BOFArgsBuffer) AddInt(d uint32) error {
	return binary.Write(b.Buffer, binary.LittleEndian, &d)
}

func (b *BOFArgsBuffer) AddString(d string) error {
	stringLen := uint32(len(d)) + 1
	err := binary.Write(b.Buffer, binary.LittleEndian, &stringLen)
	if err != nil {
		return err
	}
	dBytes := append([]byte(d), 0x00)
	return binary.Write(b.Buffer, binary.LittleEndian, dBytes)
}

func (b *BOFArgsBuffer) AddWString(d string) error {
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	strBytes := append([]byte(d), 0x00)
	utf16Data, err := encoder.Bytes(strBytes)
	if err != nil {
		return err
	}
	stringLen := uint32(len(utf16Data))
	err = binary.Write(b.Buffer, binary.LittleEndian, &stringLen)
	if err != nil {
		return err
	}
	return binary.Write(b.Buffer, binary.LittleEndian, utf16Data)
}

func (b *BOFArgsBuffer) GetBuffer() ([]byte, error) {
	outBuffer := new(bytes.Buffer)
	err := binary.Write(outBuffer, binary.LittleEndian, uint32(b.Buffer.Len()))
	if err != nil {
		return nil, err
	}
	err = binary.Write(outBuffer, binary.LittleEndian, b.Buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return outBuffer.Bytes(), nil
}
