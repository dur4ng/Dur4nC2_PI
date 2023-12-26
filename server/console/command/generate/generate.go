package generate

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/domain/models"
	serverGenerate "Dur4nC2/server/generate"
	"fmt"
	"github.com/desertbit/grumble"
	"net/url"
	"strings"
)

func parseCompileFlags(ctx *grumble.Context, con *console.ServerConsoleClient) *models.ImplantConfig {
	if ctx.Flags.String("implant") == "" {
		con.PrintErrorf("Require implant package path")
		return nil
	}
	implantPackagePath := ctx.Flags.String("implant")

	var name string
	if ctx.Flags.String("name") != "" {
		name = strings.ToLower(ctx.Flags.String("name"))
	}

	httpC2, err := ParseHTTPc2(ctx.Flags.String("http"))
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}

	os := ctx.Flags.String("os")

	isSharedLib := false
	isService := false
	isShellcode := false
	format := ctx.Flags.String("format")
	//runAtLoad := false
	var configFormat models.OutputFormat
	switch format {
	case "exe":
		configFormat = models.OutputFormat_EXECUTABLE
	case "shared":
		configFormat = models.OutputFormat_SHARED_LIB
		isSharedLib = true
		//runAtLoad = ctx.Flags.Bool("run-at-load")
	case "shellcode":
		configFormat = models.OutputFormat_SHELLCODE
		isShellcode = true
	case "service":
		configFormat = models.OutputFormat_SERVICE
		isService = true
	default:
		// Default to exe
		configFormat = models.OutputFormat_EXECUTABLE
	}

	config := &models.ImplantConfig{
		Name:               name,
		Format:             configFormat,
		IsService:          isService,
		IsShellcode:        isShellcode,
		IsSharedLib:        isSharedLib,
		URL:                httpC2.String(),
		Domain:             httpC2.Host,
		PathPrefix:         serverGenerate.Emtpy,
		ImplantPackagePath: implantPackagePath,
		OS:                 os,
	}

	return config
}

func ParseHTTPc2(arg string) (*url.URL, error) {
	arg = strings.ToLower(arg)
	var uri *url.URL
	var err error
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		uri, err = url.Parse(arg)
		if err != nil {
			return nil, err
		}
	} else {
		uri, err = url.Parse(fmt.Sprintf("https://%s", arg))
		if err != nil {
			return nil, err
		}
	}
	uri.Path = strings.TrimSuffix(uri.Path, "/")
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return nil, fmt.Errorf("invalid http(s) scheme: %s", uri.Scheme)
	}
	return uri, nil
}
