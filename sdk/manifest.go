package sdk

// Manifest represents the static metadata and security requirements for a plugin.
// This is typically defined in a plugin.yaml file stored alongside the binary.
type Manifest struct {
	// Name is the unique name of the plugin.
	Name string `yaml:"name"`
	// Version is the plugin's semantic version.
	Version string `yaml:"version"`
	// APIVersion is the compatible ZTVS API version.
	APIVersion int `yaml:"api_version"`
	// Runtime specifies the plugin execution environment.
	Runtime RuntimeMetadata `yaml:"runtime"`
	// Capabilities are the security permissions required by the plugin.
	Capabilities []string `yaml:"capabilities"`
	// Checksum is the SHA256 hex string for binary integrity.
	Checksum string `yaml:"checksum,omitempty"`
	// ChecksSupported is a list of IDs for the checks this plugin provides.
	ChecksSupported []string `yaml:"checks_supported,omitempty"`
}

// RuntimeMetadata specifies the execution environment for a plugin.
type RuntimeMetadata struct {
	// Type is the runtime type (e.g., "binary", "python", "node").
	Type string `yaml:"type"`
	// Language is the development language of the plugin.
	Language string `yaml:"language"`
	// Entrypoint is the execution path context (e.g., bin path, main.js).
	Entrypoint string `yaml:"entrypoint"`
}
