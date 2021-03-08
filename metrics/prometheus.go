package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

// InitializeMetrics incvokes custom metric functions
func InitializeMetrics() {
	fmt.Println("Init metrics")
	cpuMetrics()
	exampleGauge()
}

func exampleGauge() {
	opsQueued := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "our_company",
		Subsystem: "blob_storage",
		Name:      "ops_queued",
		Help:      "Number of blob storage operations waiting to be processed.",
	})
	prometheus.MustRegister(opsQueued)

	// 10 operations queued by the goroutine managing incoming requests.
	opsQueued.Add(10)
	// A worker goroutine has picked up a waiting operation.
	opsQueued.Dec()
	// And once more...
	opsQueued.Dec()
}

func cpuMetrics() {
	cpuGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "group_l",
		Subsystem: "minitwut",
		Name:      "cpu_usage",
		Help:      "The current CPU usage for the system",
	})

	prometheus.MustRegister(cpuGauge)

	go func() {
		for {
			percent, _ := cpu.Percent(0, true)
			cpuGauge.Set(percent[cpu.CPUser])
			time.Sleep(5 * time.Second)
			fmt.Println(percent[cpu.CPUser])
		}
	}()
}
