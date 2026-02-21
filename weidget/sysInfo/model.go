package sysinfo

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemSpec struct {
	RAM       string
	PROCESSOR string
	GPU       string
	Storage   string
}

type SysInfoWidget struct {
	Username      string
	Distro        string
	KernelVersion string
	Shell         string
	SystemSpec    SystemSpec
}

func NewSysInfoWidget() SysInfoWidget {
	S := SysInfoWidget{}
	S.GetSystemSpec()
	S.GetSysInfo()
	return S
}

func (S SysInfoWidget) Init() tea.Cmd {
	return nil
}

func (S SysInfoWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return S, nil
}

func (S SysInfoWidget) View() string {
	s := lipgloss.NewStyle().Margin(1, 2).Render(fmt.Sprintf("CPU: %s\nGPU: %s\nRAM: %s\nStorage: %s\n\n%s@%s\nKernel: %s\nShell: %s", S.SystemSpec.PROCESSOR, S.SystemSpec.GPU, S.SystemSpec.RAM, S.SystemSpec.Storage, S.Username, S.Distro, S.KernelVersion, S.Shell))
	return s
}

func (S *SysInfoWidget) GetSystemSpec() {
	// RAM
	vmStat, _ := mem.VirtualMemory()
	ram := fmt.Sprintf("%d MB", vmStat.Total/1024/1024)

	// CPU
	cpuInfo, _ := cpu.Info()
	processor := cpuInfo[0].ModelName

	// Storage
	diskStat, _ := disk.Usage("/")
	storage := fmt.Sprintf("%d GB", diskStat.Total/1024/1024/1024)

	S.SystemSpec = SystemSpec{
		RAM:       ram,
		PROCESSOR: processor,
		Storage:   storage,
		GPU:       GetGPU(), // gopsutil doesn't have GPU
	}
}

func (S *SysInfoWidget) GetSysInfo() {
	hostInfo, _ := host.Info()
	S.Username = os.Getenv("USER")
	S.Distro = hostInfo.Platform + " " + hostInfo.PlatformVersion
	S.KernelVersion = hostInfo.KernelVersion
	S.Shell = os.Getenv("SHELL")
}

func GetGPU() string {
	// Try nvidia first
	entries, err := os.ReadDir("/proc/driver/nvidia/gpus")
	if err == nil && len(entries) > 0 {
		data, _ := os.ReadFile("/proc/driver/nvidia/gpus/" + entries[0].Name() + "/information")
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "Model:") {
				return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			}
		}
	}

	// Fallback: read from /sys
	data, err := os.ReadFile("/sys/class/drm/card0/device/uevent")
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "DRIVER=") {
				return strings.SplitN(line, "=", 2)[1]
			}
		}
	}

	return "unknown"
}

func (S SysInfoWidget) SetPosition(x, y int) {
}
