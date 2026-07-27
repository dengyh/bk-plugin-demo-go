package data_api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webctx "github.com/beego/beego/v2/server/web/context"
	log "github.com/sirupsen/logrus"
)

func callTenantIDAPI(t *testing.T, tenantID string) (TenantIDResponse, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/bk_plugin/plugin_api/tenant_id", nil)
	if tenantID != "" {
		req.Header.Set(HeaderTenantID, tenantID)
	}
	req.Header.Set(HeaderRequestID, "request-tenant-001")
	recorder := httptest.NewRecorder()
	ctx := webctx.NewContext()
	ctx.Reset(recorder, req)

	var logs bytes.Buffer
	originalOutput := log.StandardLogger().Out
	originalFormatter := log.StandardLogger().Formatter
	log.SetOutput(&logs)
	log.SetFormatter(&log.JSONFormatter{DisableTimestamp: true})
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		log.SetFormatter(originalFormatter)
	})

	controller := &TenantIDController{}
	controller.Init(ctx, "TenantIDController", "Get", nil)
	controller.Get()

	var response TenantIDResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v; body=%q", err, recorder.Body.String())
	}
	return response, logs.String()
}

func TestTenantIDAPIReadsGatewayTenantHeader(t *testing.T) {
	response, logs := callTenantIDAPI(t, "tenant-blue")

	if !response.Result {
		t.Fatalf("result = false, message = %q", response.Message)
	}
	if len(response.Data) != 1 {
		t.Fatalf("data length = %d, want 1", len(response.Data))
	}
	if response.Data[0].Value != "tenant-blue" {
		t.Fatalf("tenant value = %q, want tenant-blue", response.Data[0].Value)
	}
	if response.Data[0].Label != "当前租户：tenant-blue" {
		t.Fatalf("tenant label = %q, want 当前租户：tenant-blue", response.Data[0].Label)
	}
	if !strings.Contains(logs, `"tenant_id":"tenant-blue"`) {
		t.Fatalf("log missing tenant_id: %q", logs)
	}
	if !strings.Contains(logs, `"request_id":"request-tenant-001"`) {
		t.Fatalf("log missing request_id: %q", logs)
	}
	if !strings.Contains(logs, `"msg":"data api tenant context received"`) {
		t.Fatalf("log missing success message: %q", logs)
	}
}

func TestTenantIDAPIRejectsMissingGatewayTenantHeader(t *testing.T) {
	response, logs := callTenantIDAPI(t, "")

	if response.Result {
		t.Fatal("result = true, want false")
	}
	if response.Message != "X-Bkapi-Tenant-Id header is required" {
		t.Fatalf("message = %q", response.Message)
	}
	if response.Data != nil {
		t.Fatalf("data = %#v, want nil", response.Data)
	}
	if !strings.Contains(logs, `"msg":"data api tenant context missing"`) {
		t.Fatalf("log missing error message: %q", logs)
	}
}
