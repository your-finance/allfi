// Package provider Prometheus 监控指标
package provider

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Provider 调用次数
	providerCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exchange_rate_provider_calls_total",
			Help: "Provider 调用总次数",
		},
		[]string{"provider", "status"}, // status: success, failure
	)

	// Provider 响应时间
	providerResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "exchange_rate_provider_response_seconds",
			Help:    "Provider 响应时间（秒）",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
		[]string{"provider"},
	)

	// Provider 健康状态
	providerHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "exchange_rate_provider_health",
			Help: "Provider 健康状态（1=健康，0=不健康）",
		},
		[]string{"provider"},
	)
)
