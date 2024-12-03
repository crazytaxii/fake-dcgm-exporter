package dcgm

import "fmt"

type NvidiaGPUModel string

const (
	ModelRTX4090 NvidiaGPUModel = "RTX4090"
	ModelA100    NvidiaGPUModel = "A100"
	ModelA800    NvidiaGPUModel = "A800"
)

type GPUStandard struct {
	DriverVersion          string
	ModelName              string
	FrameBufferMemoryTotal float64 // DCGM_FI_DEV_FB_USED + DCGM_FI_DEV_FB_FREE
	FreeTemperature        float64 // temperature in free status
}

var GPUMap = map[NvidiaGPUModel]GPUStandard{
	ModelRTX4090: {
		DriverVersion:          "535.183.01",
		ModelName:              "NVIDIA GeForce RTX 4090",
		FrameBufferMemoryTotal: 24216.0,
		FreeTemperature:        29.0,
	},
	ModelA100: {
		DriverVersion:          "535.104.12",
		ModelName:              "NVIDIA A100-SXM4-40GB",
		FrameBufferMemoryTotal: 40338.0,
		FreeTemperature:        29.0,
	},
	ModelA800: {
		DriverVersion:          "535.161.08",
		ModelName:              "NVIDIA A800-SXM4-80GB",
		FrameBufferMemoryTotal: 81228.0,
		FreeTemperature:        36.0,
	},
	// Add more GPUs here
}

func ErrUnknowGPUModel(model NvidiaGPUModel) error {
	return fmt.Errorf("unknown GPU model: %s", model)
}
