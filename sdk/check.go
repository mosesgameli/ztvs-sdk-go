package sdk

import "context"

// Check is the interface that all plugin checks must implement.
type Check interface {
	// ID returns a unique identifier for the check (e.g., "cve-2024-1234").
	ID() string
	// Name returns a human-readable name for the check.
	Name() string
	// Run executes the check logic.
	// It returns a Finding if a vulnerability is detected, or nil if the check passes.
	Run(ctx context.Context) (*Finding, error)
}
