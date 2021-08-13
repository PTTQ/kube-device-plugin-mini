package plugin

import (
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"os"
)

// KubeMessenger is used to communicate with api server
type KubeMessenger struct {
	client   *kubesys.KubernetesClient
	nodeName string
}

func NewKubeMessenger(masterUrl, token string) *KubeMessenger {
	client := kubesys.NewKubernetesClient(masterUrl, token)
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		log.Fatalln("Must set env NODE_NAME.")
	}
	return &KubeMessenger{
		client:   client,
		nodeName: nodeName,
	}
}
