package main

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type apiGatewayResourceSpec struct {
	Paths map[string]map[string]apiGatewayOperation `yaml:"paths"`
}

type apiGatewayOperation struct {
	OperationID string `yaml:"operationId"`
	Resource    struct {
		MatchSubpath bool `yaml:"matchSubpath"`
		Backend      struct {
			Method       string `yaml:"method"`
			Path         string `yaml:"path"`
			MatchSubpath bool   `yaml:"matchSubpath"`
		} `yaml:"backend"`
	} `yaml:"x-bk-apigateway-resource"`
}

func TestAPIGatewayResourcesExposePluginDataAPIs(t *testing.T) {
	content, err := os.ReadFile("data/api-resources.yml")
	if err != nil {
		t.Fatalf("read API Gateway resources: %v", err)
	}

	var spec apiGatewayResourceSpec
	if err := yaml.Unmarshal(content, &spec); err != nil {
		t.Fatalf("decode API Gateway resources: %v", err)
	}

	tests := []struct {
		name                string
		path                string
		method              string
		operationID         string
		backendMethod       string
		backendPath         string
		requireMatchSubpath bool
	}{
		{
			name:          "standard ops dispatch",
			path:          "/plugin_api_dispatch",
			method:        "post",
			operationID:   "plugin_api_dispatch",
			backendMethod: "post",
			backendPath:   "/{env.api_sub_path}bk_plugin/plugin_api_dispatch/",
		},
		{
			name:                "direct plugin data API",
			path:                "/bk_plugin/plugin_api/",
			method:              "x-bk-apigateway-method-any",
			operationID:         "plugin_api",
			backendMethod:       "any",
			backendPath:         "/{env.api_sub_path}bk_plugin/plugin_api/",
			requireMatchSubpath: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			methods, ok := spec.Paths[tt.path]
			if !ok {
				t.Fatalf("API Gateway path %q is not declared", tt.path)
			}
			operation, ok := methods[tt.method]
			if !ok {
				t.Fatalf("API Gateway path %q has no method %q", tt.path, tt.method)
			}
			if operation.OperationID != tt.operationID {
				t.Fatalf("operationId = %q, want %q", operation.OperationID, tt.operationID)
			}
			if operation.Resource.Backend.Method != tt.backendMethod {
				t.Fatalf("backend method = %q, want %q", operation.Resource.Backend.Method, tt.backendMethod)
			}
			if operation.Resource.Backend.Path != tt.backendPath {
				t.Fatalf("backend path = %q, want %q", operation.Resource.Backend.Path, tt.backendPath)
			}
			if tt.requireMatchSubpath &&
				(!operation.Resource.MatchSubpath || !operation.Resource.Backend.MatchSubpath) {
				t.Fatal("plugin data API resource must match subpaths")
			}
		})
	}
}
