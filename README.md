# ZTVS Go SDK

A dedicated Go SDK for building security plugins for the **Zero Trust Vulnerability Scanner (ZTVS)**. This SDK abstracts the underlying JSON-RPC communication and provides a high-level API for defining and running security checks.

## Installation

```bash
go get github.com/mosesgameli/ztvs-sdk-go
```

## Quick Start

Building a ZTVS plugin in Go is straightforward. You define one or more checks and pass them to the SDK's `Run` function.

### 1. Define a Check

Implement the `sdk.Check` interface:

```go
type MyCheck struct{}

func (c *MyCheck) ID() string   { return "check-id" }
func (c *MyCheck) Name() string { return "My Custom Check" }

func (c *MyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
    // Perform your security logic here
    if vulnFound {
        return &sdk.Finding{
            ID:          "finding-01",
            Severity:    "high",
            Title:       "Vulnerability Detected",
            Description: "A more detailed description of the flaw.",
            Evidence:    map[string]interface{}{"path": "/usr/bin/unsafe"},
            Remediation: "Patch the system immediately.",
        }, nil
    }
    return nil, nil // Return nil finding if the check passes
}
```

### 2. Run the Plugin

Bootstrap your plugin in `main.go`:

```go
package main

import (
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

func main() {
	meta := sdk.Metadata{
		Name:       "my-custom-plugin",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	checks := []sdk.Check{
		&MyCheck{},
	}

	sdk.Run(meta, checks)
}
```

## Core Concepts

### `Check` Interface
The primary contract for a plugin check.
- `ID()`: A unique identifier for the check (e.g., `cve-2024-1234`).
- `Run()`: The execution logic. Returns a `Finding` if an issue is found, or `nil` if the check passes.

### `Finding` Struct
Represents a security vulnerability.
- `Severity`: One of `critical`, `high`, `medium`, `low`, or `info`.
- `Evidence`: A map of arbitrary data to provide proof of the finding.

## Development

If you are developing the SDK locally alongside the main ZTVS repository, you can use a `replace` directive in your `go.mod`:

```go
replace github.com/mosesgameli/ztvs-sdk-go => ../sdk/go
```

## License
MIT
