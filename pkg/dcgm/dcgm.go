package dcgm

const (
	defaultDriverVersion = "535.104.12"
	defaultModelName     = "NVIDIA A100-SXM4-40GB"
	defaultGPUsNumber    = 8
)

type Workload struct {
	PodName       string `yaml:"podName,omitempty"`
	Namespace     string `yaml:"namespace,omitempty"`
	ContainerName string `yaml:"containerName,omitempty"`
	GPUs          uint32 `yaml:"gpus,omitempty"`
}

type GPUInfo struct {
	Hostname      string
	ExporterPod   string // exporter Pod name
	ModelName     string `yaml:"modelName,omitempty"`
	DriverVersion string `yaml:"driverVersion,omitempty"`
	Number        uint32 `yaml:"number,omitempty"` // number of GPUs: default 8
}

type FakeGPUConfig struct {
	*GPUInfo  `yaml:",inline"`
	Workloads []*Workload `yaml:"workloads,omitempty"`
}

// exporterPod is the exporter Pod name
func DefaultGPUConfig(hostname, exporterPod string) *FakeGPUConfig {
	return &FakeGPUConfig{
		GPUInfo: &GPUInfo{
			Hostname:      hostname,
			ExporterPod:   exporterPod,
			ModelName:     defaultModelName,
			DriverVersion: defaultDriverVersion,
			Number:        defaultGPUsNumber,
		},
	}
}
