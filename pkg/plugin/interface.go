package plugin

// DevicePlugin interface
type DevicePlugin interface {
	// Start the device plugin
	Start() error
	// Stop the device plugin and cleanup
	Stop() error
}
