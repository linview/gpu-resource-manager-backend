package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"mock-tinkerbell/config"
	"mock-tinkerbell/models"
	"github.com/sirupsen/logrus"
)

type HardwareService struct {
	hardware map[string]*models.Hardware
	mutex    sync.RWMutex
	logger   *logrus.Logger
	config   *config.Config
}

func NewHardwareService(cfg *config.Config) *HardwareService {
	service := &HardwareService{
		hardware: make(map[string]*models.Hardware),
		logger:   logrus.New(),
		config:   cfg,
	}
	service.initializeDefaultHardware()
	return service
}

func (service *HardwareService) initializeDefaultHardware() {
	servers := []struct {
		name  string
		model string
	}{
		{"gpu-server-01", "Dell PowerEdge R740"},
		{"gpu-server-02", "Dell PowerEdge R740"},
		{"gpu-server-03", "HPE ProLiant DL380"},
		{"gpu-server-04", "HPE ProLiant DL380"},
		{"gpu-server-05", "Supermicro SYS-4029GP-TRT"},
	}

	for i, serverInfo := range servers {
		server := models.NewServer(serverInfo.name, serverInfo.model)
		server.IPAddress = fmt.Sprintf("192.168.1.%d", 100+i+1)
		server.MACAddress = fmt.Sprintf("00:15:5d:01:ca:%02x", i+1)
		server.SerialNumber = fmt.Sprintf("SRV%06d", i+1)
		service.addGPUsToServer(server)
		service.hardware[server.ID] = server
	}
}

func (service *HardwareService) addGPUsToServer(server *models.Hardware) {
	gpuSpecs := service.config.Hardware.DefaultGPUSpecs
	if len(gpuSpecs) == 0 {
		return
	}
	
	numGPUTypes := rand.Intn(3) + 1
	for i := 0; i < numGPUTypes; i++ {
		spec := gpuSpecs[rand.Intn(len(gpuSpecs))]
		count := spec.Count
		if count <= 0 {
			count = rand.Intn(4) + 1
		}
		
		for j := 0; j < count; j++ {
			gpu := models.GPUSpec{
				Model:         spec.Model,
				MemoryGB:      spec.MemoryGB,
				DriverVersion: spec.DriverVersion,
				CUDAVersion:   spec.CUDAVersion,
				SerialNumber:  fmt.Sprintf("GPU%08d", rand.Intn(99999999)),
				Temperature:   rand.Intn(30) + 40,
				PowerUsage:    rand.Intn(200) + 100,
				Utilization:   rand.Intn(100),
			}
			server.AddGPU(gpu)
		}
	}
}

func (service *HardwareService) GetAllHardware() []*models.Hardware {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	hardware := make([]*models.Hardware, 0, len(service.hardware))
	for _, hw := range service.hardware {
		hardware = append(hardware, hw)
	}
	return hardware
}

func (service *HardwareService) GetHardwareByID(id string) (*models.Hardware, bool) {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	hw, exists := service.hardware[id]
	return hw, exists
}

func (service *HardwareService) UpdateHardwareStatus(id, status string) error {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	hw, exists := service.hardware[id]
	if !exists {
		return fmt.Errorf("hardware not found: %s", id)
	}
	
	hw.UpdateStatus(status)
	return nil
}

func (service *HardwareService) UpdateGPUStatus(hardwareID, gpuID, status string) error {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	hw, exists := service.hardware[hardwareID]
	if !exists {
		return fmt.Errorf("hardware not found: %s", hardwareID)
	}
	
	hw.UpdateGPUStatus(gpuID, status)
	return nil
}

func (service *HardwareService) SimulateHardwareDiscovery() {
	time.Sleep(time.Duration(service.config.Hardware.MockDelay) * time.Millisecond)
	
	service.mutex.Lock()
	defer service.mutex.Unlock()
	
	for _, hw := range service.hardware {
		if hw.Status == "unknown" {
			hw.UpdateStatus("discovered")
		}
	}
}

func (service *HardwareService) GetHardwareStats() map[string]interface{} {
	service.mutex.RLock()
	defer service.mutex.RUnlock()
	
	stats := map[string]interface{}{
		"total_servers": 0,
		"total_gpus":    0,
		"status_counts": map[string]int{},
		"gpu_models":    map[string]int{},
	}
	
	for _, hw := range service.hardware {
		if hw.Type == "server" {
			stats["total_servers"] = stats["total_servers"].(int) + 1
			stats["total_gpus"] = stats["total_gpus"].(int) + len(hw.Specs.GPU)
			
			for _, gpu := range hw.Specs.GPU {
				count := stats["gpu_models"].(map[string]int)[gpu.Model]
				stats["gpu_models"].(map[string]int)[gpu.Model] = count + 1
			}
		}
		
		count := stats["status_counts"].(map[string]int)[hw.Status]
		stats["status_counts"].(map[string]int)[hw.Status] = count + 1
	}
	
	return stats
}
