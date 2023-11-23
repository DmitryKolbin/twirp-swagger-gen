package main

import (
	"errors"
	"flag"

	"github.com/DmitryKolbin/twirp-swagger-gen/internal/swagger"
	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/protobuf/compiler/protogen"
)

var _ = spew.Dump

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	var flags flag.FlagSet
	hostname := flags.String("hostname", "example.com", "")
	pathPrefix := flags.String("path_prefix", "/twirp", "")
	outputSuffix := flags.String("output_suffix", ".swagger.json", "")
	authMode := flags.String("auth_mode", "", "bearer|custom-header\tsupport bearer (via swagger 2 API key) and api key in custom header")
	customHeader := flags.String("custom_header", "", "custom header for 'custom-header' auth_mode")
	version := flags.String("version", "", "your api version")

	opts := protogen.Options{
		ParamFunc: flags.Set,
	}
	opts.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			in := f.Desc.Path()
			log.Debugf("generating: %q", in)

			if !f.Generate {
				log.Debugf("skip generating: %q", in)
				continue
			}

			var swaggerOpts []swagger.SwaggerOpt

			if authMode != nil && *authMode != "" {
				switch *authMode {
				case "bearer":
					swaggerOpts = append(swaggerOpts, swagger.WithBearerAuthentication())
				case "custom-header":
					if customHeader == nil || *customHeader == "" {
						log.Warnf("missing custom header for: %q", *authMode)
					} else {
						swaggerOpts = append(swaggerOpts, swagger.WithApiKeyAuthentication(*customHeader))
					}
				default:
					log.Warnf("unsupported auth mode: %q", *authMode)
				}
			}
			if version != nil && *version != "" {
				swaggerOpts = append(swaggerOpts, swagger.WithVersion(*version))
			}
			writer := swagger.NewWriter(in, *hostname, *pathPrefix, swaggerOpts...)

			if err := writer.WalkFile(); err != nil {
				if errors.Is(err, swagger.ErrNoServiceDefinition) {
					log.Debugf("skip writing file, %s: %q", err, in)
					continue
				}
				return err
			}

			out := f.GeneratedFilenamePrefix + *outputSuffix
			g := gen.NewGeneratedFile(out, f.GoImportPath)
			if _, err := g.Write(writer.Get()); err != nil {
				return err
			}
		}
		return nil
	})
}
