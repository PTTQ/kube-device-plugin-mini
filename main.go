package main

import (
	"flag"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"kube-device-plugin-mini/pkg/plugin/common"
	"kube-device-plugin-mini/pkg/plugin/nvidia"
	"syscall"
)

const (
	defaultMasterUrl = "https://133.133.135.42:6443"
	defaultToken     = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN2U3VZUk16R3ZfZGNaMkw4bVktVGlRWnJGZFB2NWprU1lrd0hObnNBVFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi10Z202ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI0MjNlMDJmLTdmYzAtNDEzYi04ODczLTc0YTM3MTFkMzdkOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.KVJ7NC4NAWViLy2YkFFzzg0G4NcKnAZzw8VYooyXaLQlyfJWysR0giU8QLcSRs5BqIagff2EcVBuVHmSE4o1Zt3AMayStk-stwdtQre28adKYwR4aJLtfa1Wqmw--RiBHZmOjOmzynDdtWEe_sJPl4bGSxMvjFEKy6OepXOctnqZjUq4x2mMK-FID5hmeoHY6oAcfrRuAJsHRuLEAJQzLiMAf9heTuRNxcv3OTyfGtLOOj9risr59wilC_JWVPC5DC5TkEe4-8OeWg_mKA-lwSss_nyGMCsBqPIdPeyd3RQQ9ADPDq-JP2Nci0zoqOEwgZu3nQ3wOovR7lFBbRxsQQ"
)

var (
	masterUrl = flag.String("masterUrl", defaultMasterUrl, "Kubernetes master url.")
	token     = flag.String("token", defaultToken, "Kubernetes client token.")
)

func main() {
	flag.Parse()

	log.Infoln("Loading NVML...")
	if err := nvml.Init(); err != nil {
		log.Warningf("Failed to initialize NVML: %s.", err)
		log.Warningln("If this is a GPU node, did you set the docker default runtime to `nvidia`?")
		log.Fatalln("Kube Device Plugin fail.")
	}
	defer func() {
		log.Infof("Shutdown of NVML returned: %s.", nvml.Shutdown())
	}()

	sigChan := common.NewOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	watcher, err := common.NewFileWatcher(pluginapi.DevicePluginPath)
	if err != nil {
		log.Fatalf("Failed to created file watcher: %s.", err)
	}

	devicePlugin := nvidia.NewNvidiaDevicePlugin(*masterUrl, *token)

	go func() {
		select {
		case sig := <-sigChan:
			devicePlugin.Stop()
			log.Fatalf("Received signal %v, shutting down.", sig)
		}
	}()

restart:

	devicePlugin.Stop()
	if err := devicePlugin.Start(); err != nil {
		log.Warningf("Device plugin failed to start due to %s.", err)
		goto restart
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				log.Infof("Inotify: %s created, restarting.", pluginapi.KubeletSocket)
				goto restart
			}
		case err := <-watcher.Errors:
			log.Warningf("Inotify: %s", err)
		}
	}

}
