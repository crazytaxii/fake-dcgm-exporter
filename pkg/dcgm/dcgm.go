package dcgm

const (
	defaultGPUsNumber = 8
)

type Workload struct {
	PodName       string `yaml:"podName,omitempty"`
	Namespace     string `yaml:"namespace,omitempty"`
	ContainerName string `yaml:"containerName,omitempty"`
	GPUs          uint32 `yaml:"gpus,omitempty"`
}

type GPUInfo struct {
	Model    NvidiaGPUModel `yaml:"model,omitempty"`
	Hostname string
	Number   uint32 `yaml:"number,omitempty"` // number of GPUs: default 8
}

type FakeGPUConfig struct {
	*GPUInfo  `yaml:",inline"`
	Workloads []*Workload `yaml:"workloads,omitempty"`
}

// exporterPod is the exporter Pod name
func DefaultGPUConfig(hostname string) *FakeGPUConfig {
	return &FakeGPUConfig{
		GPUInfo: &GPUInfo{
			Hostname: hostname,
			Number:   defaultGPUsNumber,
			Model:    ModelA100, // NVIDIA A100-SXM4-40GB by default
		},
	}
}
