package main

import (
	"reflect"
	"testing"

	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
)

func TestRegisterPluginsIncludesTenantValidationVersion(t *testing.T) {
	registerPlugins()

	got := hub.GetPluginVersions()
	want := []string{"1.0.1", "1.0.0"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin versions = %v, want %v", got, want)
	}

	detail, err := hub.GetPluginDetail("1.0.1")
	if err != nil {
		t.Fatalf("GetPluginDetail(1.0.1) error = %v", err)
	}
	properties, ok := detail.ContextInputsSchemaJSON()["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("context input properties type = %T, want map[string]interface{}", detail.ContextInputsSchemaJSON()["properties"])
	}
	if _, ok := properties["tenant_id"]; !ok {
		t.Fatalf("context input properties = %v, want tenant_id", properties)
	}
}
