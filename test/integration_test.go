package test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

// mockCheck implements sdk.Check for integration testing.
type mockCheck struct {
	id string
}

func (m *mockCheck) ID() string   { return m.id }
func (m *mockCheck) Name() string { return "Integration Mock Check" }
func (m *mockCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	return &sdk.Finding{
		ID:       "int-finding-01",
		Severity: "high",
		Title:    "Integration Finding",
	}, nil
}

// TestSDKIntegration validates the end-to-end JSON-RPC flow at the SDK level.
func TestSDKIntegration(t *testing.T) {
	meta := sdk.Metadata{
		Name:       "int-test-plugin",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	checks := []sdk.Check{
		&mockCheck{id: "check-1"},
	}

	// 1. Test Handshake
	handshakeReq := rpc.Request{
		JSONRPC: "2.0",
		ID:      "h1",
		Method:  "handshake",
	}
	reqBytes, _ := json.Marshal(handshakeReq)
	in := bytes.NewReader(reqBytes)
	out := &bytes.Buffer{}

	// We use the exported Run but it's hard to test because it uses os.Stdin/Stdout.
	// Since we've refactored for internal 'run' (r, w), we can test 'run' directly,
	// but standard projects often provide a exported variant for testing or
	// set Stdin/Stdout. Here we test the newly refactored internal 'run'
	// to verify the logic.

	sdk.InternalRun(in, out, meta, checks)

	var handshakeResp rpc.Response[rpc.HandshakeResponse]
	if err := json.Unmarshal(out.Bytes(), &handshakeResp); err != nil {
		t.Fatalf("failed to decode handshake response: %v", err)
	}

	if handshakeResp.Result.Name != "int-test-plugin" {
		t.Errorf("expected plugin name int-test-plugin, got %s", handshakeResp.Result.Name)
	}

	// 2. Test Run Check
	runReq := rpc.Request{
		JSONRPC: "2.0",
		ID:      "r1",
		Method:  "run_check",
		Params:  rpc.RunCheckRequest{CheckID: "check-1"},
	}
	reqBytes, _ = json.Marshal(runReq)
	in = bytes.NewReader(reqBytes)
	out.Reset()

	sdk.InternalRun(in, out, meta, checks)

	var runResp rpc.Response[rpc.RunCheckResponse]
	if err := json.Unmarshal(out.Bytes(), &runResp); err != nil {
		t.Fatalf("failed to decode run_check response: %v", err)
	}

	if runResp.Result.Status != "fail" { // Because mockCheck returns a finding
		t.Errorf("expected status fail, got %s", runResp.Result.Status)
	}
	if runResp.Result.Finding.ID != "int-finding-01" {
		t.Errorf("expected finding ID int-finding-01, got %s", runResp.Result.Finding.ID)
	}
}
