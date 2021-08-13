package nvidia

import (
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"kube-device-plugin-mini/pkg/plugin"
)

const (
	serverSock = pluginapi.DevicePluginPath + "kubegpushare.sock"
	pluginName = "Kube-GPU-Sharing"
)

// NvidiaDevicePlugin implements the Kubernetes device plugin API
type NvidiaDevicePlugin struct {
	devices           []*pluginapi.Device
	deviceNameByIndex map[uint]string
	server            *grpc.Server
	socket            string
	stop              chan struct{}
	health            chan *pluginapi.Device
	messenger         *plugin.KubeMessenger
}

func NewNvidiaDevicePlugin(masterUrl, token string) *NvidiaDevicePlugin {
	devices, deviceMap := getDevices()
	messenger := plugin.NewKubeMessenger(masterUrl, token)

	return &NvidiaDevicePlugin{
		devices:           devices,
		deviceNameByIndex: deviceMap,
		server:            nil,
		socket:            serverSock,
		stop:              make(chan struct{}),
		health:            make(chan *pluginapi.Device),
		messenger:         messenger,
	}
}

func (p *NvidiaDevicePlugin) Name() string {
	return pluginName
}

func (p *NvidiaDevicePlugin) Start() error {
	return nil
}

func (p *NvidiaDevicePlugin) Stop() error {
	return nil
}
