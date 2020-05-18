/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package collection provides functions used by the health manager to collect
// health related metrics for FeG services and the system
package collection

import (
	"fmt"
	"time"

	"magma/feg/cloud/go/protos"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// CollectSystemStats fills out the SystemHealthStats protos with system related
// metrics that are useful to the health manager. If a metric cannot be retrieved,
// an error is returned along with all other collected metrics
func CollectSystemStats() (*protos.SystemHealthStats, error) {
	stats := &protos.SystemHealthStats{
		Time: uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}
	cpuUtilPctArray, cpuErr := cpu.Percent(0, false)
	if cpuErr == nil && len(cpuUtilPctArray) == 1 {
		stats.CpuUtilPct = float32(cpuUtilPctArray[0]) / 100
	}
	virtualMem, vmErr := mem.VirtualMemory()
	if vmErr == nil {
		stats.MemTotalBytes = virtualMem.Total
		stats.MemAvailableBytes = virtualMem.Available
	}
	if cpuErr != nil || vmErr != nil {
		return stats, fmt.Errorf("Error collecting system stats; CPU Result: %v, MEM Result: %v,", cpuErr, vmErr)
	}
	return stats, nil
}