package common

import (
	"encoding/json"
	"errors"
	"github.com/kubesys/kubernetes-client-go/pkg/kubesys"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	. "kube-device-plugin-mini/pkg/constant"
	"os"
)

// KubeMessenger is used to communicate with api server
type KubeMessenger struct {
	client   *kubesys.KubernetesClient
	nodeName string
}

func NewKubeMessenger(masterUrl, token string) *KubeMessenger {
	client := kubesys.NewKubernetesClient(masterUrl, token)
	client.Init()
	nodeName, err := os.Hostname()
	if err != nil || nodeName == "" {
		log.Fatalln("Failed to get node name.")
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

func (m *KubeMessenger) PatchGPUCount(count, core uint) error {
	node := m.getNode()
	if node == nil {
		log.Warningln("Failed to get node.")
		return errors.New("node not found")
	}

	newNode := node.DeepCopy()
	newNode.Status.Capacity[ResourceCount] = *resource.NewQuantity(int64(count), resource.DecimalSI)
	newNode.Status.Allocatable[ResourceCount] = *resource.NewQuantity(int64(count), resource.DecimalSI)
	newNode.Status.Capacity[ResourceCore] = *resource.NewQuantity(int64(core), resource.DecimalSI)
	newNode.Status.Allocatable[ResourceCore] = *resource.NewQuantity(int64(core), resource.DecimalSI)

	err := m.updateNodeStatus(newNode)
	if err != nil {
		log.Warningln("Failed to update Capacity gpu-count %s and core %s.", ResourceCount, ResourceCore)
		return errors.New("patch node status fail")
	}

	return nil
}

func (m* KubeMessenger) updateNodeStatus(node *v1.Node) error {
	nodeJson, err := json.Marshal(node)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateResourceStatus(string(nodeJson))
	if err != nil {
		return err
	}
	return nil
}

func (m *KubeMessenger) GetPendingPodsOnNode() []v1.Pod {
	var pods []v1.Pod
	podList, _ := m.client.ListResources("Pod", "")
	podListBytes, _ := json.Marshal(podList.Object)
	var podListObject v1.PodList
	json.Unmarshal(podListBytes, &podListObject)
	for _, pod := range podListObject.Items {
		if pod.Spec.NodeName == m.nodeName && pod.Status.Phase == "Pending" {
			pod.Kind = "Pod"
			pod.APIVersion = "v1"
			pods = append(pods, pod)
		}
	}
	return pods
}

func (m *KubeMessenger) UpdatePodAnnotations(pod *v1.Pod) error {
	podJson, err := json.Marshal(pod)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateResource(string(podJson))
	if err != nil {
		return err
	}
	return nil
}

func (m *KubeMessenger) GetPodOnNode(podName, namespace string) *v1.Pod {
	pod, err := m.client.GetResource("Pod", namespace, podName)
	if err != nil {
		return nil
	}
	podBytes, _ := json.Marshal(pod.Object)
	var out v1.Pod
	err = json.Unmarshal(podBytes, &out)
	if err != nil {
		return nil
	}
	return &out
}