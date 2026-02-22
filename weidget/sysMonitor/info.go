package sysmonitor

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemStats struct {
	CPUPercent  float64
	RAMPercent  float64
	RAMUsed     uint64
	RAMTotal    uint64
	DiskPercent float64
	// DiskRead    uint64 // bytes/s
	// DiskWrite   uint64 // bytes/s
	// GPUPercent  float64
	// GPUMemUsed  uint64
	logger string
	mu     sync.RWMutex
}

func (s *SystemStats) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// CPU — blocks 500ms internally, so this loop naturally
			// runs at ~500ms intervals
			cpuPct, _ := cpu.Percent(500*time.Millisecond, false)

			ram, _ := mem.VirtualMemory()
			du, _ := disk.Usage("/")

			s.mu.Lock()
			if len(cpuPct) > 0 {
				s.CPUPercent = cpuPct[0]
			}
			s.RAMPercent = ram.UsedPercent
			s.RAMUsed = ram.Used
			s.RAMTotal = ram.Total
			s.DiskPercent = du.UsedPercent
			s.mu.Unlock()
		}
	}()
}

func (s *SystemStats) Read() SystemStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s // copy — caller gets a snapshot
}
