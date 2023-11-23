package swagger

import "github.com/go-openapi/spec"

func WithBearerAuthentication() func(target *spec.Swagger) {
	return func(target *spec.Swagger) {
		if target.SecurityDefinitions == nil {
			target.SecurityDefinitions = make(map[string]*spec.SecurityScheme)
		}

		target.SecurityDefinitions["bearer"] = &spec.SecurityScheme{
			VendorExtensible: spec.VendorExtensible{},
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Description: "Enter `Bearer: {token}`",
				Type:        "apiKey",
				Name:        "Authorization",
				In:          "header",
			},
		}

		if target.Security == nil {
			target.Security = make([]map[string][]string, 0)
		}

		target.Security = append(target.Security, map[string][]string{"bearer": {}})
	}
}

func WithApiKeyAuthentication(customApiHeader string) func(target *spec.Swagger) {
	return func(target *spec.Swagger) {
		if target.SecurityDefinitions == nil {
			target.SecurityDefinitions = make(map[string]*spec.SecurityScheme)
		}

		target.SecurityDefinitions["api-key"] = &spec.SecurityScheme{
			VendorExtensible: spec.VendorExtensible{},
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type: "apiKey",
				Name: customApiHeader,
				In:   "header",
			},
		}

		if target.Security == nil {
			target.Security = make([]map[string][]string, 0)
		}

		target.Security = append(target.Security, map[string][]string{"api-key": {}})
	}
}

func WithVersion(version string) func(target *spec.Swagger) {
	return func(target *spec.Swagger) {
		if target.Info == nil {
			target.Info = &spec.Info{}
		}
		target.Info.Version = version
	}
}
