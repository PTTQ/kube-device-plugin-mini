package constant

import pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

const (
	DefaultMasterUrl = "https://133.133.135.42:6443"
	DefaultToken     = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN2U3VZUk16R3ZfZGNaMkw4bVktVGlRWnJGZFB2NWprU1lrd0hObnNBVFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi10Z202ZyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjI0MjNlMDJmLTdmYzAtNDEzYi04ODczLTc0YTM3MTFkMzdkOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.KVJ7NC4NAWViLy2YkFFzzg0G4NcKnAZzw8VYooyXaLQlyfJWysR0giU8QLcSRs5BqIagff2EcVBuVHmSE4o1Zt3AMayStk-stwdtQre28adKYwR4aJLtfa1Wqmw--RiBHZmOjOmzynDdtWEe_sJPl4bGSxMvjFEKy6OepXOctnqZjUq4x2mMK-FID5hmeoHY6oAcfrRuAJsHRuLEAJQzLiMAf9heTuRNxcv3OTyfGtLOOj9risr59wilC_JWVPC5DC5TkEe4-8OeWg_mKA-lwSss_nyGMCsBqPIdPeyd3RQQ9ADPDq-JP2Nci0zoqOEwgZu3nQ3wOovR7lFBbRxsQQ"

	ServerSock    = pluginapi.DevicePluginPath + "doslab.sock"
	ResourceName  = "doslab.io/gpu-memory"
	ResourceCount = "doslab.io/gpu-count"
	ResourceCore  = "doslab.io/gpu-core"

	AnnResourceAssumeTime = "doslab.io/gpu-assume-time"
	AnnGemSchedulerIp     = "doslab.io/gem-scheduler-ip"
	AnnAssignedFlag       = "doslab.io/gpu-assigned"
	AnnResourceUUID       = "doslab.io/gpu-uuid"
	AnnGemPodManagerPort  = "doslab.io/gem-podmanager-port"

	EnvResourceUUID            = "DOSLAB_IO_GPU_UUID"
	EnvResourceUsedByPod       = "DOSLAB_IO_GPU_RESOURCE_USED_BY_POD"
	EnvResourceUsedByContainer = "DOSLAB_IO_GPU_RESOURCE_USED_BY_CONTAINER"
	EnvResourceTotal           = "DOSLAB_IO_GPU_RESOURCE_TOTAL"

	EnvNvidiaGPU                = "NVIDIA_VISIBLE_DEVICES"
	EnvPodName                  = "POD_NAME"
	EnvPodManagerPort           = "POD_MANAGER_PORT"
	EnvPodManagerIp             = "POD_MANAGER_IP"
	EnvLDPreload                = "LD_PRELOAD"
	EnvNvidiaDriverCapabilities = "NVIDIA_DRIVER_CAPABILITIES"

	KubeShareLibraryPath = "/kubeshare/library"
)
