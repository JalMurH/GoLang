package gui

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	Pid       int32  `json:"pid"`
	Name      string `json:"name"`
	CreateTms int64  `json:"create_time"`
}

func GetParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})

	// Obtener información del CPU
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	params["cpu"] = cpuInfo

	// Obtener información de la memoria
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	params["memory"] = memInfo

	// Obtener información del disco
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	params["disk"] = diskInfo

	// Obtener información de las particiones
	diskPInfo, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	params["partitions"] = diskPInfo

	// Obtener información de las tarjetas gráficas
	// gpuInfo, err := getGPUinfo()
	// if err != nil {
	// 	return nil, err
	// }
	// params["gpu"] = gpuInfo

	// Obtener información de los procesos en ejecución
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	var processList []ProcessInfo
	for _, p := range processes {
		name, _ := p.Name()
		createTms, _ := p.CreateTime()
		processList = append(processList, ProcessInfo{p.Pid, name, createTms})
	}

	var recentProcesses []ProcessInfo
	for i := 0; i < 10 && i < len(processList); i++ {
		recentProcesses = append(recentProcesses, processList[i])
	}
	params["processes"] = recentProcesses

	//obtener datos de temperatura
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}
	params["sensors"] = sensors

	return params, nil
}
