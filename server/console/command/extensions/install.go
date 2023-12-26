package extensions

import (
	"Dur4nC2/server/console"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desertbit/grumble"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var installedExtensions = map[string]*ExtensionManifest{}

type ExtensionManifest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Help        string               `json:"help"`
	Files       []*extensionFile     `json:"files"`
	Arguments   []*extensionArgument `json:"arguments"`
	Entrypoint  string               `json:"entrypoint"`
	Init        string               `json:"init"`
	DependsOn   string               `json:"depends_on"`
}

type extensionFile struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Path string `json:"path"`
}

type extensionArgument struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Desc     string `json:"desc"`
	Optional bool   `json:"optional"`
}

func ExtensionInstallerCmd(ctx *grumble.Context, con *console.ServerConsoleClient) (string, error) {
	dirPath := ctx.Args.String("dir-path")
	extensionsManifests, err := loadExtendionsManifest(dirPath)
	if err != nil {
		return "", err
	}
	installExtensions(extensionsManifests)
	con.PrintSuccessf("New extension installed")
	return "", nil
}

func loadExtendionsManifest(manifestPath string) (*ExtensionManifest, error) {
	data, err := ioutil.ReadFile(manifestPath)
	//fmt.Println(data)
	if err != nil {
		return nil, err
	}
	ext, err := parseExtensionManifest(data)
	if err != nil {
		return nil, err
	}
	return ext, nil
}
func parseExtensionManifest(data []byte) (*ExtensionManifest, error) {
	extManifest := &ExtensionManifest{}
	err := json.Unmarshal(data, &extManifest)
	if err != nil {
		return nil, err
	}
	if extManifest.Name == "" {
		return nil, errors.New("missing `name` field in extension manifest")
	}
	if len(extManifest.Files) == 0 {
		return nil, errors.New("missing `files` field in extension manifest")
	}
	for _, extFiles := range extManifest.Files {
		if extFiles.OS == "" {
			return nil, errors.New("missing `files.os` field in extension manifest")
		}
		if extFiles.Arch == "" {
			return nil, errors.New("missing `files.arch` field in extension manifest")
		}
		extFiles.Path = ResolvePath(extFiles.Path)
		if extFiles.Path == "" || extFiles.Path == "/" {
			return nil, errors.New("missing `files.path` field in extension manifest")
		}
		extFiles.OS = strings.ToLower(extFiles.OS)
		extFiles.Arch = strings.ToLower(extFiles.Arch)
	}
	if extManifest.Help == "" {
		return nil, errors.New("missing `help` field in extension manifest")
	}
	return extManifest, nil
}
func installExtensions(manifests *ExtensionManifest) {
	installedExtensions[manifests.Name] = manifests
}
func ResolvePath(in string) string {
	if strings.Contains(in, ":\\") {
		parts := strings.Split(in, ":\\")
		in = parts[len(parts)-1]
	}
	out := filepath.Clean(fmt.Sprintf("x:\\%s", in))
	return strings.TrimPrefix(out, "x:")
}
