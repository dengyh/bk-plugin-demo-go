package v101

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

func newTestContext(t *testing.T, tenantID string) (*kit.Context, *bytes.Buffer, *objectStore) {
	t.Helper()

	inputs, err := json.Marshal(Inputs{Hello: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	contextInputs, err := json.Marshal(ContextInputs{Executor: "tester", TenantID: tenantID})
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
	entry := log.NewEntry(logger).WithField("trace_id", "trace-tenant-001")

	return kit.NewContext(
		"trace-tenant-001",
		constants.StateEmpty,
		1,
		reader,
		contextStore,
		outputsStore,
		entry,
	), &logs, outputsStore
}

func TestExecuteLogsTenantIDWithTraceID(t *testing.T) {
	ctx, logs, outputsStore := newTestContext(t, "tenant-blue")

	if err := (&Plugin{}).Execute(ctx); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	outputs, ok := outputsStore.value.(*Outputs)
	if !ok {
		t.Fatalf("outputs type = %T, want *Outputs", outputsStore.value)
	}
	if outputs.World != "hello" {
		t.Fatalf("World = %q, want %q", outputs.World, "hello")
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(logs.Bytes()), &entry); err != nil {
		t.Fatalf("decode log entry: %v; log=%q", err, logs.String())
	}
	if entry["trace_id"] != "trace-tenant-001" {
		t.Fatalf("trace_id = %v, want trace-tenant-001", entry["trace_id"])
	}
	if entry["tenant_id"] != "tenant-blue" {
		t.Fatalf("tenant_id = %v, want tenant-blue", entry["tenant_id"])
	}
	if entry["msg"] != "tenant context received" {
		t.Fatalf("msg = %v, want tenant context received", entry["msg"])
	}
}

func TestExecuteRejectsMissingTenantID(t *testing.T) {
	ctx, logs, outputsStore := newTestContext(t, "")

	err := (&Plugin{}).Execute(ctx)
	if err == nil {
		t.Fatal("Execute() error = nil, want missing tenant_id error")
	}
	if !strings.Contains(err.Error(), "tenant_id is required") {
		t.Fatalf("Execute() error = %q, want tenant_id is required", err)
	}
	if outputsStore.value != nil {
		t.Fatalf("outputs = %#v, want nil", outputsStore.value)
	}
	if !strings.Contains(logs.String(), `"trace_id":"trace-tenant-001"`) {
		t.Fatalf("log missing trace_id: %q", logs.String())
	}
	if !strings.Contains(logs.String(), `"msg":"tenant context missing"`) {
		t.Fatalf("log missing tenant context error: %q", logs.String())
	}
}
