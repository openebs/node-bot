/*
Copyright 2018 The OpenEBS Author

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// StartingTime is the starting time of ndm.
	StartingTime time.Time
	// Uptime is the uptime of ndm.
	Uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ndm_uptime_seconds",
		Help: "Uptime of node disk manager.",
	})
	collectorUptimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "collector", "uptime_seconds"),
		"Uptime of collector.",
		[]string{"collector"},
		nil,
	)
)

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector() (Collector, error) {
	var diskLabelNames = []string{"device"}

	return &diskstatsCollector{
		ignoredDevicesPattern: regexp.MustCompile(ignoredDevices),
		descs: []typedFactorDesc{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"),
					"Total number of reads completed successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"),
					"Total number of bytes read successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"),
					"Total number of writes completed successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "write_bytes_total"),
					"Total number of bytes wrote successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_average"),
					"Average number of bytes read successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "write_bytes_average"),
					"Average number of bytes wrote successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "read_latency"),
					"Latency in read operations.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "write_latency"),
					"Latency in write operations.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
		},
	}, nil
}
