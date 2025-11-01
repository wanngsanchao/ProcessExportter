package main

import (
	"github.com/shirou/gopsutil/v3/process"
)

//get the mem and cpu usage of the service_test process
func GetResourceUsage(pid int) (cpuusgae,memusage float64,err error) {
    p,err :=  process.NewProcess(int32(pid))

    if err != nil {
    return 0,0,err
    }

    //cpu percent
    cpu,err := p.CPUPercent()

    if err != nil {
        return 0,0,err
    }
    //mem usage
    mem,err := p.MemoryInfo()

    if err != nil {
        return 0,0,err
    }

    return cpu,float64(mem.RSS) / 1024 / 1024,nil
}
