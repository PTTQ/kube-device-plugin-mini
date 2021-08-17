package nvidia

import (
	"context"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	log "github.com/sirupsen/logrus"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"strings"
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

func watchXIDs(ctx context.Context, devs []*pluginapi.Device, xids chan<- *pluginapi.Device) {
	eventSet := nvml.NewEventSet()
	defer nvml.DeleteEventSet(eventSet)

	for _, d := range devs {
		err := nvml.RegisterEventForDevice(eventSet, nvml.XidCriticalError, d.ID)
		if err != nil && strings.HasSuffix(err.Error(), "Not Supported") {
			log.Warningf("Warning: %s is too old to support healthchecking: %s. Marking it unhealthy.", d.ID, err)
			xids <- d
			continue
		}
		check(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		e, err := nvml.WaitForEvent(eventSet, 5000)
		if err != nil && e.Etype != nvml.XidCriticalError {
			continue
		}

		// FIXME: formalize the full list and document it.
		// http://docs.nvidia.com/deploy/xid-errors/index.html#topic_4
		// Application errors: the GPU should still be healthy
		if e.Edata == 31 || e.Edata == 43 || e.Edata == 45 {
			continue
		}

		if e.UUID == nil || len(*e.UUID) == 0 {
			// All devices are unhealthy
			for _, d := range devs {
				xids <- d
			}
			continue
		}

		for _, d := range devs {
			if d.ID == *e.UUID {
				log.Printf("XidCriticalError: Xid=%d on Device=%s, the device will go unhealthy.", e.Edata, d.ID)
				xids <- d
			}
		}
	}

}