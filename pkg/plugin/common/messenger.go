package common

import (
	"encoding/json"
	"errors"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"os"
)

const (
	resourceCount = "kubesys.io/gpu-count"
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
		log.Fatalln("There is no env NODE_NAME.")
	}
	return &KubeMessenger{
		client:   client,
		nodeName: nodeName,
	}
}

func (m *KubeMessenger) getNode() *v1.Node {
	node, err := m.client.GetResource("Node", "", m.nodeName)
	if err != nil {
		return nil
	}
	nodeBytes, _ := json.Marshal(node.Object)
	var out v1.Node
	err = json.Unmarshal(nodeBytes, &out)
	if err != nil {
		return nil
	}
	return &out
}

func (m *KubeMessenger) PatchGPUCount(gpuCount int) error {
	node := m.getNode()
	if node == nil {
		log.Warningln("Failed to get node.")
		return errors.New("node not found")
	}

	if val, ok := node.Status.Capacity[resourceCount]; ok {
		if val.Value() == int64(gpuCount) {
			log.Infof("No need to update Capacity %s", resourceCount)
			return nil
		}
	}

	newNode := node.DeepCopy()
	newNode.Status.Capacity[resourceCount] = *resource.NewQuantity(int64(gpuCount), resource.DecimalSI)
	newNode.Status.Allocatable[resourceCount] = *resource.NewQuantity(int64(gpuCount), resource.DecimalSI)

	err := m.patchNodeStatus(newNode)
	if err != nil {
		log.Warningln("Failed to update Capacity %s.", resourceCount)
		return errors.New("patch node status fail")
	}

	return nil
}

func (m* KubeMessenger) patchNodeStatus(newNode *v1.Node) error {
	nodeJson, err := json.Marshal(newNode)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateResourceStatus(string(nodeJson))
	if err != nil {
		return err
	}
}