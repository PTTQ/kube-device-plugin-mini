package nvidia

import (
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// NvidiaDevicePlugin implements the Kubernetes device plugin API
type NvidiaDevicePlugin struct {
	devices           []*pluginapi.Device
	server            *grpc.Server
	socket            string
	stop              chan struct{}
	deviceNameByIndex map[uint]string
}

func NewNvidiaDevicePlugin() *NvidiaDevicePlugin {
	
}