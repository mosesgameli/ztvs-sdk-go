package sdk

// Metadata represents the plugin information provided to the ZTVS host during handshake.
type Metadata struct {
	// Name is the name of the plugin for display and registration.
	Name string
	// Version is the current version of the plugin.
	Version string
	// APIVersion is the compatible ZTVS communication API version.
	APIVersion int
}
