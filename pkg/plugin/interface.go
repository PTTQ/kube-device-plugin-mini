package plugin

// DevicePlugin interface
type DevicePlugin interface {
	// Get the device plugin name
	Name() string
	// Start the device plugin
	Start() error
	// Stop the device plugin
	Stop() error
}
