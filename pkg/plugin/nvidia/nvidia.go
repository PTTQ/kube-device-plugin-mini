package nvidia

import (
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	log "github.com/sirupsen/logrus"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func check(err error) {
	if err != nil {
		log.Fatalf("Fatal: %s", err)
	}
}

func transformDevice(d *nvml.Device) *pluginapi.Device {
	device := pluginapi.Device{}
	log.Infof("Device %s' Path is %s.", d.UUID, d.Path)
	device.ID = d.UUID
	device.Health = pluginapi.Healthy
	if d.CPUAffinity != nil {
		device.Topology = &pluginapi.TopologyInfo{
			Nodes: []*pluginapi.NUMANode{
				&pluginapi.NUMANode{
					ID: int64(*(d.CPUAffinity)),
				},
			},
		}
	}
	return &device
}

func getDevices() ([]*pluginapi.Device, map[uint]string) {
	n, err := nvml.GetDeviceCount()
	check(err)

	log.Infoln("Getting devices...")
	var devices []*pluginapi.Device
	var deviceMap map[uint]string
	for idx := uint(0); idx < n; idx++ {
		d, err := nvml.NewDevice(idx)
		check(err)
		devices = append(devices, transformDevice(d))
		deviceMap[idx] = devices[idx].ID
	}
	return devices, deviceMap
}
