package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bk-plugin-demo-go/data_api"

	"github.com/TencentBlueKing/beego-runtime/routers"
	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
	"github.com/beego/beego/v2/server/web"
)

func TestRegisterPluginsIncludesTenantValidationVersion(t *testing.T) {
	registerPlugins()

	got := hub.GetPluginVersions()
	want := []string{"1.0.3", "1.0.2", "1.0.1", "1.0.0"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin versions = %v, want %v", got, want)
	}

	detail, err := hub.GetPluginDetail("1.0.3")
	if err != nil {
		t.Fatalf("GetPluginDetail(1.0.3) error = %v", err)
	}
	properties, ok := detail.ContextInputsSchemaJSON()["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("context input properties type = %T, want map[string]interface{}", detail.ContextInputsSchemaJSON()["properties"])
	}
	if _, ok := properties["tenant_id"]; !ok {
		t.Fatalf("context input properties = %v, want tenant_id", properties)
	}

	inputProperties := detail.InputsSchemaJSON()["properties"].(map[string]interface{})
	tenantField := inputProperties["data_api_tenant_id"].(map[string]interface{})
	component := tenantField["ui:component"].(map[string]interface{})
	props := component["props"].(map[string]interface{})
	remoteConfig := props["remoteConfig"].(map[string]interface{})
	wantURL := "{{ $context.get('bk_plugin_api_host')['bkplugin-go2'] + 'bk_plugin/plugin_api/tenant_id' }}"
	if remoteConfig["url"] != wantURL {
		t.Fatalf("remoteConfig.url = %v, want %v", remoteConfig["url"], wantURL)
	}
}

func TestTenantIDDataAPIDispatch(t *testing.T) {
	web.BConfig.CopyRequestBody = true
	web.AddNamespace(routers.PluginApiNamespace)

	body := bytes.NewBufferString(`{
		"url": "/bk_plugin/plugin_api/tenant_id",
		"method": "GET",
		"username": "tester",
		"data": {}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/bk_plugin/plugin_api_dispatch", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(data_api.HeaderTenantID, "tenant-blue")
	req.Header.Set(data_api.HeaderRequestID, "request-tenant-dispatch-001")
	recorder := httptest.NewRecorder()

	web.BeeApp.Handlers.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%q", recorder.Code, recorder.Body.String())
	}
	var response data_api.TenantIDResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v; body=%q", err, recorder.Body.String())
	}
	if !response.Result || len(response.Data) != 1 || response.Data[0].Value != "tenant-blue" {
		t.Fatalf("response = %#v", response)
	}
}
