module kube-device-plugin-mini

go 1.16

require (
	github.com/NVIDIA/gpu-monitoring-tools v0.0.0-20210811154823-a8b4da88dc2c
	github.com/fsnotify/fsnotify v1.4.9
	github.com/kubesys/kubernetes-client-go v0.7.0
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.38.0
	k8s.io/api v0.22.0
	k8s.io/kubelet v0.22.0
)
