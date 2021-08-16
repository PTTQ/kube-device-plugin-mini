package nvidia

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"kube-device-plugin-mini/pkg/plugin/common"
	"net"
	"os"
	"path"
	"time"
)

const (
	serverSock   = pluginapi.DevicePluginPath + "kubegpushare.sock"
	resourceName = "kubesys.io/gpu"
)

// NvidiaDevicePlugin implements the Kubernetes device plugin API
type NvidiaDevicePlugin struct {
	devices           []*pluginapi.Device
	deviceNameByIndex map[uint]string
	server            *grpc.Server
	socket            string
	stop              chan struct{}
	health            chan *pluginapi.Device
	messenger         *common.KubeMessenger
}

func NewNvidiaDevicePlugin(masterUrl, token string) *NvidiaDevicePlugin {
	messenger := common.NewKubeMessenger(masterUrl, token)

	return &NvidiaDevicePlugin{
		devices:           nil,
		deviceNameByIndex: nil,
		server:            nil,
		socket:            serverSock,
		stop:              nil,
		health:            nil,
		messenger:         messenger,
	}
}

// Start mainly starts the gRPC server and register the device plugin to Kubelet
func (p *NvidiaDevicePlugin) Start() error {
	log.Infoln("Starting the device plugin")

	p.devices, p.deviceNameByIndex = getDevices()
	p.server = grpc.NewServer([]grpc.ServerOption{}...)
	p.stop = make(chan struct{})
	p.health = make(chan *pluginapi.Device)

	err := p.messenger.PatchGPUCount(len(p.devices))
	if err != nil {
		p.cleanup()
		return err
	}

	err = p.Serve()
	if err != nil {
		p.cleanup()
		log.Warningln("Could not serve: %s.", err)
		return err
	}

	err = p.Register(pluginapi.KubeletSocket, resourceName)
	if err != nil {
		p.cleanup()
		log.Infof("Could not register device plugin: %s.", err)
		return err
	}
	log.Infoln("Registered device plugin with Kubelet.")

	go p.healthCheck()

	return nil
}

func (p *NvidiaDevicePlugin) Stop() error {
	log.Infoln("Stopping the device plugin")

	if p == nil || p.server == nil {
		return nil
	}
	p.server.Stop()
	if err := os.Remove(p.socket); err != nil && !os.IsNotExist(err) {
		return err
	}
	p.cleanup()
	return nil
}

func (p *NvidiaDevicePlugin) Serve() error {
	sock, err := net.Listen("unix", p.socket)
	if err != nil {
		return err
	}
	pluginapi.RegisterDevicePluginServer(p.server, p)

	go p.server.Serve(sock)

	conn, err := dial(p.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

func (p *NvidiaDevicePlugin) Register(endpoint, resourceName string) error {
	conn, err := dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	req := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(endpoint),
		ResourceName: resourceName,
	}

	_, err = client.Register(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (p *NvidiaDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {

}

func (p *NvidiaDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {

}

func (p *NvidiaDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (p *NvidiaDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (p *NvidiaDevicePlugin) GetPreferredAllocation(context.Context, *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

func (p *NvidiaDevicePlugin) cleanup() {
	close(p.stop)
	p.devices = nil
	p.deviceNameByIndex = nil
	p.server = nil
	p.stop = nil
	p.health = nil
}

func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c, err := grpc.DialContext(ctx, unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *NvidiaDevicePlugin) healthCheck() {

}
