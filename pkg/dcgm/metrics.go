package dcgm

import (
	"bytes"
	"cmp"
	"fmt"
	"text/template"

	"github.com/google/uuid"
)

type FakeGPUMetrics struct {
	UUID          string
	Hostname      string
	ModelName     string
	DriverVersion string
	PodName       string
	PodNamespace  string
	ContainerName string
}

type MetricsRenderer struct {
	gpuInfo   *GPUInfo
	workloads []*Workload
	tmpl      *template.Template
}

func NewGPUMetricsRenderer(cfg *FakeGPUConfig) (*MetricsRenderer, error) {
	tmpl, err := template.New("fake-dcgm-metrics").Parse(metricsTemplate)
	if err != nil {
		return nil, err
	}
	return &MetricsRenderer{
		gpuInfo:   cfg.GPUInfo,
		workloads: cfg.Workloads,
		tmpl:      tmpl,
	}, nil
}

func (r *MetricsRenderer) Render() ([]byte, error) {
	// gpusMetrics is the list of fake GPU metrics
	// length of the list equals to the number of GPUs
	gpusMetrics := make([]*FakeGPUMetrics, r.gpuInfo.Number)
	for i := range gpusMetrics {
		gpusMetrics[i] = &FakeGPUMetrics{
			UUID:          uuid.NewMD5(uuid.NameSpaceURL, []byte(fmt.Sprintf("%s-nvidia%d", r.gpuInfo.Hostname, i))).String(), // generate stable UUID
			Hostname:      r.gpuInfo.Hostname,
			ModelName:     r.gpuInfo.ModelName,
			DriverVersion: r.gpuInfo.DriverVersion,
		}
	}

	var used int
	for _, workload := range r.workloads {
		for range workload.GPUs {
			gpusMetrics[used].PodName = workload.PodName
			gpusMetrics[used].PodNamespace = workload.Namespace
			gpusMetrics[used].ContainerName = cmp.Or(workload.ContainerName, fmt.Sprintf("%s-c", workload.PodName))
			if used++; used >= int(r.gpuInfo.Number) {
				// all GPUs are used by Pods
				break
			}
		}
	}

	buf := &bytes.Buffer{}
	err := r.tmpl.Execute(buf, gpusMetrics)
	return buf.Bytes(), err
}

const metricsTemplate = `# HELP DCGM_FI_DEV_SM_CLOCK SM clock frequency (in MHz).
# TYPE DCGM_FI_DEV_SM_CLOCK gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_SM_CLOCK{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 210
{{- end }}
# HELP DCGM_FI_DEV_MEM_CLOCK Memory clock frequency (in MHz).
# TYPE DCGM_FI_DEV_MEM_CLOCK gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_MEM_CLOCK{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 1215
{{- end }}
# HELP DCGM_FI_DEV_MEMORY_TEMP Memory temperature (in C).
# TYPE DCGM_FI_DEV_MEMORY_TEMP gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_MEMORY_TEMP{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 29
{{- end }}
# HELP DCGM_FI_DEV_GPU_TEMP GPU temperature (in C).
# TYPE DCGM_FI_DEV_GPU_TEMP gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_GPU_TEMP{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 30
{{- end }}
# HELP DCGM_FI_DEV_POWER_USAGE Power draw (in W).
# TYPE DCGM_FI_DEV_POWER_USAGE gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_POWER_USAGE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 55.415000
{{- end }}
# HELP DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Total energy consumption since boot (in mJ).
# TYPE DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION counter
{{- range $i, $metrics := . }}
DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 1345789458
{{- end }}
# HELP DCGM_FI_DEV_PCIE_REPLAY_COUNTER Total number of PCIe retries.
# TYPE DCGM_FI_DEV_PCIE_REPLAY_COUNTER counter
{{- range $i, $metrics := . }}
DCGM_FI_DEV_PCIE_REPLAY_COUNTER{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_GPU_UTIL GPU utilization (in %).
# TYPE DCGM_FI_DEV_GPU_UTIL gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_GPU_UTIL{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_MEM_COPY_UTIL Memory utilization (in %).
# TYPE DCGM_FI_DEV_MEM_COPY_UTIL gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_MEM_COPY_UTIL{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_ENC_UTIL Encoder utilization (in %).
# TYPE DCGM_FI_DEV_ENC_UTIL gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_ENC_UTIL{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_DEC_UTIL Decoder utilization (in %).
# TYPE DCGM_FI_DEV_DEC_UTIL gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_DEC_UTIL{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_XID_ERRORS Value of the last XID error encountered.
# TYPE DCGM_FI_DEV_XID_ERRORS gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_XID_ERRORS{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_FB_FREE Framebuffer memory free (in MiB).
# TYPE DCGM_FI_DEV_FB_FREE gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_FB_FREE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 40334
{{- end }}
# HELP DCGM_FI_DEV_FB_USED Framebuffer memory used (in MiB).
# TYPE DCGM_FI_DEV_FB_USED gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_FB_USED{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 4
{{- end }}
# HELP DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL Total number of NVLink bandwidth counters for all lanes.
# TYPE DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL counter
{{- range $i, $metrics := . }}
DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_VGPU_LICENSE_STATUS vGPU License status
# TYPE DCGM_FI_DEV_VGPU_LICENSE_STATUS gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_VGPU_LICENSE_STATUS{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS Number of remapped rows for uncorrectable errors
# TYPE DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS counter
{{- range $i, $metrics := . }}
DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS Number of remapped rows for correctable errors
# TYPE DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS counter
{{- range $i, $metrics := . }}
DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_DEV_ROW_REMAP_FAILURE Whether remapping of rows has failed
# TYPE DCGM_FI_DEV_ROW_REMAP_FAILURE gauge
{{- range $i, $metrics := . }}
DCGM_FI_DEV_ROW_REMAP_FAILURE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_PROF_GR_ENGINE_ACTIVE Ratio of time the graphics engine is active (in %).
# TYPE DCGM_FI_PROF_GR_ENGINE_ACTIVE gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_GR_ENGINE_ACTIVE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0.000000
{{- end }}
# HELP DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Ratio of cycles the tensor (HMMA) pipe is active (in %).
# TYPE DCGM_FI_PROF_PIPE_TENSOR_ACTIVE gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0.000000
{{- end }}
# HELP DCGM_FI_PROF_DRAM_ACTIVE Ratio of cycles the device memory interface is active sending or receiving data (in %).
# TYPE DCGM_FI_PROF_DRAM_ACTIVE gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_DRAM_ACTIVE{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0.000000
{{- end }}
# HELP DCGM_FI_PROF_PCIE_TX_BYTES The rate of data transmitted over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_TX_BYTES gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_PCIE_TX_BYTES{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 16264
{{- end }}
# HELP DCGM_FI_PROF_PCIE_RX_BYTES The rate of data received over the PCIe bus - including both protocol headers and data payloads - in bytes per second.
# TYPE DCGM_FI_PROF_PCIE_RX_BYTES gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_PCIE_RX_BYTES{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 11605
{{- end }}
# HELP DCGM_FI_PROF_NVLINK_TX_BYTES The number of bytes of active NvLink tx (transmit) data including both header and payload.
# TYPE DCGM_FI_PROF_NVLINK_TX_BYTES gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_NVLINK_TX_BYTES{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
# HELP DCGM_FI_PROF_NVLINK_RX_BYTES The number of bytes of active NvLink rx (read) data including both header and payload.
# TYPE DCGM_FI_PROF_NVLINK_RX_BYTES gauge
{{- range $i, $metrics := . }}
DCGM_FI_PROF_NVLINK_RX_BYTES{gpu="{{ $i }}",UUID="{{ $metrics.UUID }}",device="nvidia{{ $i }}",modelName="{{ $metrics.ModelName }}",Hostname="{{ $metrics.Hostname }}",DCGM_FI_DRIVER_VERSION="{{ $metrics.DriverVersion }}",container="{{ $metrics.ContainerName }}",namespace="{{ $metrics.PodNamespace }}",pod="{{ $metrics.PodName }}"} 0
{{- end }}
`
