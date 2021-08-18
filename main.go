package main

import (
	"flag"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	. "kube-device-plugin-mini/pkg/constant"
	"kube-device-plugin-mini/pkg/plugin/common"
	"kube-device-plugin-mini/pkg/plugin/nvidia"
	"syscall"
)

var (
	masterUrl = flag.String("masterUrl", DefaultMasterUrl, "Kubernetes master url.")
	token     = flag.String("token", DefaultToken, "Kubernetes client token.")
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

	if _, m := nvidia.GetDevices(); len(m) == 0 {
		log.Warningln("There is no device, try to restart.")
		goto restart
	}

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
