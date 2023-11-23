package main

import (
	"flag"

	"github.com/DmitryKolbin/twirp-swagger-gen/internal/swagger"
	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

var _ = spew.Dump

func parse(
	hostname string,
	filename string,
	output string,
	prefix string,
	authMode string,
	customHeader string,
	version string,
) error {
	if filename == output {
		return errors.New("output file must be different than input file")
	}

	var swaggerOpts []swagger.SwaggerOpt

	if authMode != "" {
		switch authMode {
		case "bearer":
			swaggerOpts = append(swaggerOpts, swagger.WithBearerAuthentication())
		case "custom-header":
			if customHeader == "" {
				log.Warnf("missing custom header for: %q", authMode)
			} else {
				swaggerOpts = append(swaggerOpts, swagger.WithApiKeyAuthentication(customHeader))
			}
		default:
			log.Warnf("unsupported auth mode: %q", authMode)
		}
	}
	if version != "" {
		swaggerOpts = append(swaggerOpts, swagger.WithVersion(version))
	}

	writer := swagger.NewWriter(filename, hostname, prefix, swaggerOpts...)
	if err := writer.WalkFile(); err != nil {
		if !errors.Is(err, swagger.ErrNoServiceDefinition) {
			return err
		}
	}
	return writer.Save(output)
}

func main() {
	var (
		in           string
		out          string
		host         string
		pathPrefix   string
		authMode     string
		customHeader string
		version      string
	)
	flag.StringVar(&in, "in", "", "Input source .proto file")
	flag.StringVar(&out, "out", "", "Output swagger.json file")
	flag.StringVar(&host, "host", "api.example.com", "API host name")
	flag.StringVar(&pathPrefix, "pathPrefix", "/twirp", "Twrirp server path prefix")
	flag.StringVar(&authMode, "auth_mode", "", "bearer|custom-header\tsupport bearer (via swagger 2 API key) and api key in custom header")
	flag.StringVar(&customHeader, "auth_custom_header", "", "custom header for 'custom-header' auth_mode")
	flag.StringVar(&version, "version", "", "your api version")
	flag.Parse()

	if in == "" {
		log.Fatalf("Missing parameter: -in [input.proto]")
	}
	if out == "" {
		log.Fatalf("Missing parameter: -out [output.proto]")
	}
	if host == "" {
		log.Fatalf("Missing parameter: -host [api.example.com]")
	}

	if err := parse(
		host,
		in,
		out,
		pathPrefix,
		authMode,
		customHeader,
		version,
	); err != nil {
		log.WithError(err).Fatal("exit with error")
	}
}
