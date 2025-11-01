package main

import (
    "fmt"
    "strconv"
	"github.com/prometheus/client_golang/prometheus"
)

//step 1. define the custome metrics

type ProcessInfo struct {
    Procname string
    processUpDesc *prometheus.Desc
    processCPUDesc *prometheus.Desc
    ProcessMemDesc *prometheus.Desc
}

func (p *ProcessInfo) Init() {
    //1. define the custome metric,the status of the process
    p.processUpDesc = prometheus.NewDesc(
        fmt.Sprintf("%s_up",p.Procname),
        fmt.Sprintf("whether or not the %s process is running",p.Procname),
        []string{"pid"},
        prometheus.Labels{"app":"monitor"},
    )

    //2. define the custome metric,the cpu usage of the process
    p.processCPUDesc = prometheus.NewDesc(
        fmt.Sprintf("%s_cpu_usage_percent",p.Procname),
        fmt.Sprintf("cpu usage percentage of the %s process",p.Procname),
        []string{"pid"},
        prometheus.Labels{"app":"monitor"},
    )

    //3. define the custome metric,the mem usage of the process
    p.ProcessMemDesc = prometheus.NewDesc(
        fmt.Sprintf("%s_mem_usage_percent",p.Procname),
        fmt.Sprintf("memory usage(MB) of the %s process",p.Procname),
        []string{"pid"},
        prometheus.Labels{"app":"monitor"},
    )

}


//1.register the metric to defualt regitrra
//implement the Desc func
func (p *ProcessInfo) Describe(ch chan <- *prometheus.Desc) {
    ch <- p.processUpDesc
    ch <- p.processCPUDesc
    ch <- p.ProcessMemDesc
}


//2.the func will be excuted while geting the http://127.0.0.1:8080:/metrics
//implement the Collect func
func (p *ProcessInfo) Collect(ch chan <- prometheus.Metric) {
    pids,err := GetserviceTestPids(p.Procname)

    if err != nil {
        ch <- prometheus.MustNewConstMetric(
            p.processUpDesc,
            prometheus.GaugeValue,
            0,
            "",
        )
        return
    }

    for _,pid := range pids {
        //1. get the cpu and mem usage from GetResourceUsage
        cpuusgae,memusage,err := GetResourceUsage(pid)

        if err != nil {
            ch <- prometheus.MustNewConstMetric(
                p.processUpDesc,
                prometheus.GaugeValue,
                0,
                strconv.Itoa(pid),
            )
            continue
        }
        //2.write the status of process to the ch
        ch <- prometheus.MustNewConstMetric(
            p.processUpDesc,
            prometheus.GaugeValue,
            1,
            strconv.Itoa(pid),
        )
        //3.writhe the cpu usage to the ch
        ch <- prometheus.MustNewConstMetric(
            p.processCPUDesc,
            prometheus.GaugeValue,
            cpuusgae,
            strconv.Itoa(pid),
        )

        //4. write the mem usage to the ch
        ch <- prometheus.MustNewConstMetric(
            p.ProcessMemDesc,
            prometheus.GaugeValue,
            memusage,
            strconv.Itoa(pid),
        )
    }
    
}


func InitAllProcessMetric(allprocess []string) []prometheus.Collector {
    var all []prometheus.Collector

    for _,pname := range allprocess {
        p := &ProcessInfo{}
        p.Procname = pname
        p.Init()
        all = append(all,p)
    }

    return all
}
