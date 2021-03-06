package repository

import (
	"testing"

	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricPool(t *testing.T) {
	tmp := core.NewMetricsPool()
	test := []struct {
		name string
		want *core.MetricsPool
	}{
		{
			name: "Test #1",
			want: tmp,
		},
	}
	for _, tt := range test {
		mp := core.NewMetricsPool()
		assert.Equal(t, tt.want.M, mp.M)
	}
}

// func TestNewMetrics(t *testing.T) {

// 	test := []struct {
// 		name   string
// 		values []string
// 		want   Metric
// 	}{
// 		{
// 			name:   "Test #1",
// 			values: []string{"testMetric", "counter", "1"},
// 			want:   Metric{ID: "testMetric", Delta: 1, MType: "counter"},
// 		},
// 		{
// 			name:   "Test #2",
// 			values: []string{"metric2", "gauge", "123"},
// 			want:   Metric{ID: "metric2", Value: 123, MType: "gauge"},
// 		},
// 		{
// 			name:   "Test #3",
// 			values: []string{"metrics3", "counter", "111"},
// 			want:   Metric{ID: "metrics3", Delta: 111, MType: "counter"},
// 		},
// 	}
// 	for _, tt := range test {
// 		m, status := NewMetric(tt.values[1], tt.values[0], tt.values[2])
// 		if status != "" {
// 			fmt.Println(status)
// 		}
// 		assert.Equal(t, tt.want, *m)
// 	}
// }

// func TestConvert(t *testing.T) {
// 	tmp := MetricsPool{M: map[string]Metric{}}
// 	test := []struct {
// 		name string

// 		want MetricsPool
// 	}{
// 		{
// 			name: "Test #1",
// 			want: tmp,
// 		},
// 	}
// 	for _, tt := range test {
// 		mp := NewMetricsPool()
// 		assert.Equal(t, tt.want, *mp)
// 	}
// }
