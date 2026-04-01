package main

import (
	"context"
	"fmt"

	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

// GreetCheck is a basic example check that always passes.
type GreetCheck struct{}

func (c *GreetCheck) ID() string   { return "hello-world-check" }
func (c *GreetCheck) Name() string { return "Hello World Greet Check" }

func (c *GreetCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	fmt.Println("Running GreetCheck...")
	// In a real check, you'd perform some audit logic here.
	// For this example, we'll always pass.
	return nil, nil // Passing check
}

// FailingCheck is a basic example check that always fails with a finding.
type FailingCheck struct{}

func (c *FailingCheck) ID() string   { return "failing-check" }
func (c *FailingCheck) Name() string { return "Example Failing Check" }

func (c *FailingCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	return &sdk.Finding{
		ID:          "finding-fail-ex",
		Severity:    "info",
		Title:       "Example Finding",
		Description: "This is a demonstration of a failing check.",
		Evidence: map[string]interface{}{
			"reason": "This check was designed to fail for demonstration purposes.",
		},
		Remediation: "No action required, this is just an example.",
	}, nil
}

func main() {
	// 1. Define plugin metadata
	meta := sdk.Metadata{
		Name:       "example-hello-world",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	// 2. Define the checks to run
	checks := []sdk.Check{
		&GreetCheck{},
		&FailingCheck{},
	}

	// 3. Hand off execution to the SDK's Run function
	sdk.Run(meta, checks)
}
