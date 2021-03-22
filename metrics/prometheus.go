package metrics

import (
	"fmt"
	"time"

	"github.com/heyjoakim/devops-21/services"
	"github.com/prometheus/client_golang/prometheus"
	cpu "github.com/shirou/gopsutil/cpu"
	mem "github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

var HistogramVecs map[string]*prometheus.HistogramVec
var GaugeOpts map[string]*prometheus.GaugeOpts

func GetHistogramVec(name string) (*prometheus.HistogramVec, error) {
	result := HistogramVecs[name]
	if result == nil {
		return nil, fmt.Errorf("Could not find historgram with name: %s", name)
	}
	return result, nil
}

// InitializeMetrics invokes custom metric functions
func InitializeMetrics() {
	log.Info("Init metrics")
	cpuMetrics()
	memoryMetrics()
	userCountMetrics()
	messageCountMetrics()
	apiEndpointDurationsMetrics()
}

const measurementDelay = 5

func messageCountMetrics() {
	messagesGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "group_l",
		Subsystem: "minitwut",
		Name:      "message_count",
		Help:      "The current number of messages in the database",
	})

	prometheus.MustRegister(messagesGauge)

	go func() {
		for {
			count := services.GetMessageCount()
			messagesGauge.Set(float64(count))
			time.Sleep(measurementDelay * time.Second)
		}
	}()
}

func userCountMetrics() {
	userGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "group_l",
		Subsystem: "minitwut",
		Name:      "user_count",
		Help:      "The current number of registered users in the database",
	})

	prometheus.MustRegister(userGauge)

	go func() {
		for {
			count := services.GetUserCount()
			userGauge.Set(float64(count))
			time.Sleep(measurementDelay * time.Second)
		}
	}()
}

func memoryMetrics() {
	memoryGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "group_l",
		Subsystem: "minitwut",
		Name:      "memory_usage",
		Help:      "The current memory usage for the system",
	})

	prometheus.MustRegister(memoryGauge)

	go func() {
		for {
			v, _ := mem.VirtualMemory()
			memoryGauge.Set(v.UsedPercent)
			time.Sleep(measurementDelay * time.Second)
		}
	}()
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
			cpuGauge.Set(percent[0])
			time.Sleep(measurementDelay * time.Second)
		}
	}()
}

func apiEndpointDurationsMetrics() {

}
