package sdk

// Finding represents a vulnerability or security issue discovered by a check.
type Finding struct {
	// ID is the unique identifier for this specific finding instance.
	ID string
	// Severity is the risk level (e.g., "high", "medium", "low").
	Severity string
	// Title is a short summary of the finding.
	Title string
	// Description is a detailed explanation of the vulnerability.
	Description string
	// Evidence is a map of supporting data or proof for the finding.
	Evidence map[string]interface{}
	// Remediation provides advice on how to fix the issue.
	Remediation string
}
