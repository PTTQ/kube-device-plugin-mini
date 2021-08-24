# kube-device-plugin-mini

## 项目实现

本项目[1] 依据 Kubernetes 的官方文档[2] 进行设计。在代码实现部分，参考了 volcano-sh/devices [3] 以及 AliyunContainerService/gpushare-device-plugin [4] 的实现。

## 项目结构

本项目主要由 Go 语言实现，采用 go module 方式管理依赖包，其主要目录如下所示。

```
kube-device-plugin-mini
├─pkg
│   ├─constant
│   │   └─constant.go
│   └─plugin
│       ├─common
│       │   ├─messenger.go
│       │   └─watcher.go
│       ├─nvidia
│       │   ├─nvidia.go
│       │   └─server.go
│       └─interface.go
├─go.mod
├─gpu-pod-demo.yaml
├─main.go
└─README.md
```

其中 `main.go` 是代码的入口文件， `gpu-pod-demo.yaml` 是测试用例文件， `constant.go` 文件保存了项目的一些常量， `messenger.go` 通过封装研究小组自研的 kubernetes-client-go 库[5] 来进行一些与 Kubernetes API Server 通信的操作， `watcher.go` 监听文件的增减以及操作系统的信号， `nvidia.go` 用于获取设备信息， `server.go` 维护整个插件的生命周期以及实现重要接口， `interface.go` 定义了插件的抽象接口。

## 插件的生命周期

插件以守护进程或 DaemonSet 的形式运行。

在进行一些必要的初始化工作后，插件启动，进入一个无限循环。在该循环中插件主要监测两个信号，一个是操作系统的信号 SIGHUP， SIGINT 和 SIGTERM ，一旦监听到此类信号，立刻停止并退出插件；另一个是监测主机路径  `/var/lib/kubelet/device-plugins/` 下的 `kubelet.sock` 文件，当 kubelet 重启的时候，新的 kubelet 实例会删除已经存在的 Unix 套接字，插件必须重新注册自己。

插件启动之后，首先需要获取设备列表，由于 Kubernetes 的设备资源必须为整数，而我们想要分配的资源粒度为显存 MB ，因此将物理设备抽象为显存大小数量的虚拟设备，以 `doslab.io/gpu-memory` 为资源名注册。获取完设备列表后，将 GPU 的数量 `doslab.io/gpu-count` 以及算力（一个物理 GPU 对应100算力） `doslab.io/gpu-core` 也上报到 kubelet 中，供上层调度组件使用。

获取设备列表完毕后，启动 gRPC 服务，然后向 kubelet 注册，注册需要上报设备插件的 Unix 套接字、设备插件的 API 版本以及资源名（本项目为 `doslab.io/gpu-memory` ）。

## 接口实现

插件的 gRPC 服务需要我们实现 `ListAndWatch` 和 `Allocate` 接口。

接口 `ListAndWatch` 的主要功能是返回设备列表构成的数据流，当设备发生变化或者消失时，该接口会返回新的设备列表。本项目没有实现设备的健康监测功能，默认设备启动后就一直可用，因此仅仅在调用时返回设备列表，在插件关闭时返回空的设备列表。

接口 `Allocate` 的主要功能是让容器启动并获取一些必要的信息，以便能够使用设备。插件首先获取处于 Pending 状态的 Pod 列表，过滤掉不符合条件的 Pod ，并按照 Pod 的 `DOSLAB_IO_GPU_ASSUME_TIME` 排序。然后比对资源请求数量，找到真正调用 `Allocate` 接口的 Pod 。最后将 `NVIDIA_VISIBLE_DEVICES` 等环境变量注入到容器之中，使容器能够使用对应的 GPU 设备。

## 相关地址

[1] https://github.com/PTTQ/kube-device-plugin-mini

[2] https://kubernetes.io/zh/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/

[3] https://github.com/volcano-sh/devices

[4] https://github.com/AliyunContainerService/gpushare-device-plugin

[5] https://github.com/kubesys/kubernetes-client-go