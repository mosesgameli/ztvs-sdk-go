package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
)

type mockCheck struct {
	id      string
	name    string
	finding *Finding
	err     error
}

func (m *mockCheck) ID() string   { return m.id }
func (m *mockCheck) Name() string { return m.name }
func (m *mockCheck) Run(ctx context.Context) (*Finding, error) {
	return m.finding, m.err
}

func TestRun_Handshake(t *testing.T) {
	meta := Metadata{Name: "test-plugin", Version: "1.0.0", APIVersion: 1}
	checks := []Check{&mockCheck{id: "test-check"}}

	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "handshake",
	}
	reqData, _ := json.Marshal(req)

	in := bytes.NewReader(reqData)
	var out bytes.Buffer

	InternalRun(in, &out, meta, checks)

	var resp rpc.Response[rpc.HandshakeResponse]
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.ID != "1" {
		t.Errorf("expected ID 1, got %s", resp.ID)
	}
	if resp.Result.Name != "test-plugin" {
		t.Errorf("expected plugin name test-plugin, got %s", resp.Result.Name)
	}
	if len(resp.Result.ChecksSupported) != 1 || resp.Result.ChecksSupported[0] != "test-check" {
		t.Errorf("unexpected checks supported: %v", resp.Result.ChecksSupported)
	}
}

func TestRun_RunCheck_Pass(t *testing.T) {
	meta := Metadata{Name: "test-plugin"}
	checks := []Check{&mockCheck{id: "test-check", finding: nil}}

	params := rpc.RunCheckRequest{CheckID: "test-check"}
	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "2",
		Method:  "run_check",
		Params:  params,
	}
	reqData, _ := json.Marshal(req)

	in := bytes.NewReader(reqData)
	var out bytes.Buffer

	InternalRun(in, &out, meta, checks)

	var resp rpc.Response[rpc.RunCheckResponse]
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Result.Status != "pass" {
		t.Errorf("expected status pass, got %s", resp.Result.Status)
	}
	if resp.Result.Finding != nil {
		t.Error("expected nil finding for passing check")
	}
}

func TestRun_RunCheck_Fail(t *testing.T) {
	meta := Metadata{Name: "test-plugin"}
	finding := &Finding{
		ID:       "f1",
		Severity: "high",
		Title:    "Vulnerability",
	}
	checks := []Check{&mockCheck{id: "test-check", finding: finding}}

	params := rpc.RunCheckRequest{CheckID: "test-check"}
	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "3",
		Method:  "run_check",
		Params:  params,
	}
	reqData, _ := json.Marshal(req)

	in := bytes.NewReader(reqData)
	var out bytes.Buffer

	InternalRun(in, &out, meta, checks)

	var resp rpc.Response[rpc.RunCheckResponse]
	if err := json.Unmarshal(out.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Result.Status != "fail" {
		t.Errorf("expected status fail, got %s", resp.Result.Status)
	}
	if resp.Result.Finding == nil || resp.Result.Finding.ID != "f1" {
		t.Errorf("unexpected finding in response: %v", resp.Result.Finding)
	}
}

func TestRun_UnknownMethod(t *testing.T) {
	meta := Metadata{}
	checks := []Check{}

	req := rpc.Request{
		JSONRPC: "2.0",
		ID:      "4",
		Method:  "unknown",
	}
	reqData, _ := json.Marshal(req)

	in := bytes.NewReader(reqData)
	var out bytes.Buffer

	InternalRun(in, &out, meta, checks)

	var resp map[string]interface{}
	json.Unmarshal(out.Bytes(), &resp)

	errObj := resp["error"].(map[string]interface{})
	if errObj["code"].(float64) != -32601 {
		t.Errorf("expected error code -32601, got %v", errObj["code"])
	}
}
