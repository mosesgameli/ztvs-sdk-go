package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
)

// Run is the main entry point for a ZTVS plugin.
// It reads JSON-RPC requests from standard input and writes responses to standard output.
func Run(meta Metadata, checks []Check) {
	InternalRun(os.Stdin, os.Stdout, meta, checks)
}

// InternalRun is used for integration testing. It reads from r and writes to w.
func InternalRun(r io.Reader, w io.Writer, meta Metadata, checks []Check) {
	var req rpc.Request
	err := json.NewDecoder(r).Decode(&req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode request: %v\n", err)
		os.Exit(1)
	}

	switch req.Method {
	case "handshake":
		ids := make([]string, len(checks))
		for i, c := range checks {
			ids[i] = c.ID()
		}

		resp := rpc.Response[rpc.HandshakeResponse]{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: rpc.HandshakeResponse{
				Name:            meta.Name,
				Version:         meta.Version,
				APIVersion:      meta.APIVersion,
				ChecksSupported: ids,
			},
		}
		_ = json.NewEncoder(w).Encode(resp)

	case "run_check":
		// Handle params unmarshaling
		var runReq rpc.RunCheckRequest
		b, _ := json.Marshal(req.Params)
		_ = json.Unmarshal(b, &runReq)

		for _, c := range checks {
			if c.ID() == runReq.CheckID {
				finding, err := c.Run(context.Background())
				if err != nil {
					sendError(w, req.ID, 5000, err.Error())
					return
				}

				status := "pass"
				var rpcFinding *rpc.Finding
				if finding != nil {
					status = "fail"
					rpcFinding = &rpc.Finding{
						ID:          finding.ID,
						CheckID:     c.ID(),
						Severity:    finding.Severity,
						Title:       finding.Title,
						Description: finding.Description,
						Evidence:    finding.Evidence,
						Remediation: finding.Remediation,
					}
				}

				resp := rpc.Response[rpc.RunCheckResponse]{
					JSONRPC: "2.0",
					ID:      req.ID,
					Result: rpc.RunCheckResponse{
						Status:  status,
						Finding: rpcFinding,
					},
				}

				_ = json.NewEncoder(w).Encode(resp)
				return
			}
		}
		sendError(w, req.ID, 4002, "check not found")
	default:
		sendError(w, req.ID, -32601, "method not found")
	}
}

func sendError(w io.Writer, id string, code int, msg string) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": msg,
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}
