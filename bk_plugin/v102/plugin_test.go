package v102

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/TencentBlueKing/bk-plugin-framework-go/constants"
	"github.com/TencentBlueKing/bk-plugin-framework-go/kit"
	log "github.com/sirupsen/logrus"
)

type jsonContextReader struct {
	inputs        []byte
	contextInputs []byte
}

func (r *jsonContextReader) ReadInputs(v interface{}) error {
	return json.Unmarshal(r.inputs, v)
}

func (r *jsonContextReader) ReadContextInputs(v interface{}) error {
	return json.Unmarshal(r.contextInputs, v)
}

type objectStore struct {
	value interface{}
}

func (s *objectStore) Write(_ string, v interface{}) error {
	s.value = v
	return nil
}

func (s *objectStore) Read(_ string, _ interface{}) error {
	return nil
}

func newTestContext(t *testing.T, contextTenantID, dataAPITenantID string) (*kit.Context, *bytes.Buffer, *objectStore) {
	t.Helper()

	inputs, err := json.Marshal(Inputs{Hello: "hello", DataAPITenantID: dataAPITenantID})
	if err != nil {
		t.Fatal(err)
	}
	contextInputs, err := json.Marshal(ContextInputs{Executor: "tester", TenantID: contextTenantID})
	if err != nil {
		t.Fatal(err)
	}

	reader := &jsonContextReader{inputs: inputs, contextInputs: contextInputs}
	contextStore := &objectStore{}
	outputsStore := &objectStore{}
	var logs bytes.Buffer
	logger := log.New()
	logger.SetOutput(&logs)
	logger.SetFormatter(&log.JSONFormatter{DisableTimestamp: true})
	entry := log.NewEntry(logger).WithField("trace_id", "trace-tenant-002")

	return kit.NewContext(
		"trace-tenant-002",
		constants.StateEmpty,
		1,
		reader,
		contextStore,
		outputsStore,
		entry,
	), &logs, outputsStore
}

func TestInputsFormLoadsTenantIDThroughDataAPI(t *testing.T) {
	var form map[string]interface{}
	if err := json.Unmarshal(InputsFormJSON, &form); err != nil {
		t.Fatalf("decode form: %v", err)
	}

	properties := form["properties"].(map[string]interface{})
	tenantField := properties["data_api_tenant_id"].(map[string]interface{})
	component := tenantField["ui:component"].(map[string]interface{})
	props := component["props"].(map[string]interface{})
	remoteConfig := props["remoteConfig"].(map[string]interface{})
	if remoteConfig["url"] != "/plugin_service/data_api/bk-plugin-demo-go/bk_plugin/plugin_api/tenant_id" {
		t.Fatalf("remoteConfig.url = %v", remoteConfig["url"])
	}
}

func TestExecuteLogsMatchingContextAndDataAPITenantIDs(t *testing.T) {
	ctx, logs, outputsStore := newTestContext(t, "tenant-blue", "tenant-blue")

	if err := (&Plugin{}).Execute(ctx); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	outputs, ok := outputsStore.value.(*Outputs)
	if !ok {
		t.Fatalf("outputs type = %T, want *Outputs", outputsStore.value)
	}
	if outputs.World != "hello" {
		t.Fatalf("World = %q, want hello", outputs.World)
	}
	if !strings.Contains(logs.String(), `"trace_id":"trace-tenant-002"`) {
		t.Fatalf("log missing trace_id: %q", logs.String())
	}
	if !strings.Contains(logs.String(), `"tenant_id":"tenant-blue"`) {
		t.Fatalf("log missing context tenant_id: %q", logs.String())
	}
	if !strings.Contains(logs.String(), `"data_api_tenant_id":"tenant-blue"`) {
		t.Fatalf("log missing data_api_tenant_id: %q", logs.String())
	}
	if !strings.Contains(logs.String(), `"msg":"tenant context validation succeeded"`) {
		t.Fatalf("log missing success message: %q", logs.String())
	}
}

func TestExecuteRejectsTenantIDMismatch(t *testing.T) {
	ctx, logs, outputsStore := newTestContext(t, "tenant-blue", "tenant-red")

	err := (&Plugin{}).Execute(ctx)
	if err == nil {
		t.Fatal("Execute() error = nil, want tenant mismatch error")
	}
	if !strings.Contains(err.Error(), "tenant_id mismatch") {
		t.Fatalf("Execute() error = %q", err)
	}
	if outputsStore.value != nil {
		t.Fatalf("outputs = %#v, want nil", outputsStore.value)
	}
	if !strings.Contains(logs.String(), `"msg":"tenant context validation failed"`) {
		t.Fatalf("log missing mismatch message: %q", logs.String())
	}
}
