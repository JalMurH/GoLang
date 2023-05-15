package gui

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type GPUInfo struct {
	Description string `json:"description"`
	Product     string `json:"product"`
	Vendor      string `json:"vendor"`
	Temperature float64
}

func getGPUinfo() ([]GPUInfo, error) {
	cmd := exec.Command("lshw", "-C", "display", "-json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var gpuInfo []GPUInfo
	err = json.Unmarshal(out.Bytes(), &gpuInfo)
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("sensors")
	out.Reset()
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out.String(), "\n")
	re := regexp.MustCompile(`\+([\d.]+)°C`)
	for _, line := range lines {
		if strings.Contains(line, "temp1") {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				temp, err := strconv.ParseFloat(match[1], 64)
				if err == nil {
					for i := range gpuInfo {
						gpuInfo[i].Temperature = temp
					}
				}
			}
			break
		}
	}

	return gpuInfo, nil
}
