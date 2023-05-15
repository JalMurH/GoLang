package gui

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func Graph(params map[string]interface{}) {

	myApp := app.New()
	myWindow := myApp.NewWindow("Inference")

	cpuInfo := params["cpu"].([]cpu.InfoStat)
	memInfo := params["memory"].(*mem.VirtualMemoryStat)
	diskUInfo := params["disk"].(*disk.UsageStat)
	diskPInfo := params["partitions"].([]disk.PartitionStat)
	var partitionLabels []fyne.CanvasObject
	for _, partition := range diskPInfo {
		usageStat, _ := disk.Usage(partition.Mountpoint)
		label := widget.NewLabel(fmt.Sprintf("%s: %.2f GB / %.2f GB", partition.Mountpoint, float64(usageStat.Used)/1024/1024/1024, float64(usageStat.Total)/1024/1024/1024))
		partitionLabels = append(partitionLabels, label)
	}
	processes := params["processes"].([]ProcessInfo)
	var processLabels []fyne.CanvasObject
	for _, process := range processes {
		proLabel := widget.NewLabel(fmt.Sprintf("PID: %d Name: %s", process.Pid, process.Name))
		processLabels = append(processLabels, proLabel)
	}
	//sensors := params["sensors"]
	//gpuInfo := params["gpu"].([]GPUInfo)

	cpuMNStr := fmt.Sprintf("CPU Model: %v", cpuInfo[0].ModelName)
	cpuStr := fmt.Sprintf("CPU Cores: %v", cpuInfo[0].Cores)
	cpuHzStr := fmt.Sprintf("CPU Speed: %.2f GHz", float64(cpuInfo[0].Mhz)/1000)
	memStr := fmt.Sprintf("Memory: %.2f GB / %.2f GB", float64(memInfo.Used)/1024/1024/1024, float64(memInfo.Total)/1024/1024/1024)
	diskUStr := fmt.Sprintf("Disk Usage: %.2f%%", diskUInfo.UsedPercent)
	diskTStr := fmt.Sprintf("Disk Capacity: %.2f GB", float64(diskUInfo.Total)/1024/1024/1024)
	diskNStr := fmt.Sprintf("Number of Disks: %v", len(diskPInfo))
	//processesStr := fmt.Sprintf("Processes: %v", len(processes))
	//sensorsStr := fmt.Sprintf("%v", sensors)
	//gpuInfoStr := fmt.Sprintf("%v", gpuInfo)

	//fyne
	cpuMNL := widget.NewLabel(cpuMNStr)
	cpuL := widget.NewLabel(cpuStr)
	cpuHzL := widget.NewLabel(cpuHzStr)
	cpuCorL := widget.NewLabel(cpuStr)
	memL := widget.NewLabel(memStr)
	disUkL := widget.NewLabel(diskUStr)
	disTkL := widget.NewLabel(diskTStr)
	disNkL := widget.NewLabel(diskNStr)
	updateL := widget.NewLabel("Updating...")
	//processL := widget.NewLabel("Processes: " + processesStr)
	//gpuL := widget.NewLabel("GPU: " + gpuInfoStr)
	//sensorsL := widget.NewLabel("Sensors: " + sensorsStr)

	//Crear contenedores para la información de la CPU, memoria, disco y particiones
	//bgRect := canvas.NewRectangle(color.RGBA{R: 120, G: 140, B: 180, A: 255})
	cpuContainer := container.NewVBox(
		widget.NewLabelWithStyle(cpuMNL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(cpuL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(cpuHzL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(cpuCorL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	// memContainer := container.NewVBox(
	// 	widget.NewLabelWithStyle(memL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	// )
	diskContainer := container.NewVBox(
		widget.NewLabelWithStyle(disUkL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(disTkL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(disNkL.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	partitionContainer := container.NewVBox(partitionLabels...)
	processContainer := container.NewVBox(processLabels...)
	scrollableContainer := container.NewVScroll(processContainer)

	// Crear tarjetas para la información de la CPU, memoria, disco y particiones
	cpuCard := widget.NewCard("CPU", "Información de la CPU", cpuContainer)
	memCard := widget.NewCard("RAM", "Información de la memoria RAM", memL)
	diskCard := widget.NewCard("Disk", "Información del Almacenamiento", diskContainer)
	partitionsCard := widget.NewCard("Partitions", "Información de las particiones", partitionContainer)
	processCard := widget.NewCard("Process", "PID	Name", scrollableContainer)
	updateCard := widget.NewCard("Date", "GMT-5", updateL)

	colums := container.New(layout.NewHBoxLayout(), cpuCard, memCard, diskCard, partitionsCard, processCard)
	raws := container.New(layout.NewVBoxLayout(), updateCard)
	grid := container.New(layout.NewVBoxLayout(), colums, raws)

	myWindow.SetContent(grid)
	myWindow.Show()
	go func() {
		for {
			time.Sleep(time.Second)

			// Get the latest data
			params, err := GetParams()
			if err != nil {
				log.Println(err)
				continue
			}

			// Update the labels with the new data
			cpuInfo := params["cpu"].([]cpu.InfoStat)
			memInfo := params["memory"].(*mem.VirtualMemoryStat)
			diskUInfo := params["disk"].(*disk.UsageStat)

			cpuMNL.SetText(fmt.Sprintf("CPU Model: %v", cpuInfo[0].ModelName))
			cpuMNL.Refresh()
			cpuL.SetText(fmt.Sprintf("CPU Cores: %v", cpuInfo[0].Cores))
			cpuL.Refresh()
			cpuHzL.SetText(fmt.Sprintf("CPU Speed: %.2f GHz", float64(cpuInfo[0].Mhz)/1000))
			cpuHzL.Refresh()
			memL.SetText(fmt.Sprintf("Memory: %.2f GB / %.2f GB", float64(memInfo.Used)/1024/1024/1024, float64(memInfo.Total)/1024/1024/1024))
			memL.Refresh()
			disUkL.SetText(fmt.Sprintf("Disk Usage: %.2f%%", diskUInfo.UsedPercent))
			disUkL.Refresh()
			disTkL.SetText(fmt.Sprintf("Disk Capacity: %.2f GB", float64(diskUInfo.Total)/1024/1024/1024))
			disTkL.Refresh()

			// Update the update label
			updateL.SetText(fmt.Sprintf("Updated at %s", time.Now().Format(time.RFC1123)))
			updateL.Refresh()
			//Actualizar Cards

			cpuContainer.Refresh()
			diskContainer.Refresh()
			partitionContainer.Refresh()
			processContainer.Refresh()
			scrollableContainer.Refresh()

			cpuCard.Refresh()
			diskCard.Refresh()
			memCard.Refresh()
			partitionsCard.Refresh()
			processCard.Refresh()
			updateCard.Refresh()

			// Refresh the user interface
			myWindow.Canvas().Refresh(grid)
		}
	}()
	myApp.Run()
}
