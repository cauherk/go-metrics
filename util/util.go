package util

import (
	"github.com/cauherk/go-metrics"
	"github.com/cauherk/go-metrics/impl"
	"time"
)

var mr metrics.MetricRegistry

func NewOrGetMetricRegistry(maxMetricCount int) metrics.MetricRegistry {
	if mr == nil {
		mr = impl.NewMetricRegistry(maxMetricCount)
	}
	return mr
}

func NewMetricName(key string, tags map[string]string) *metrics.MetricName {
	return metrics.NewMetricName(key, tags, metrics.Major)
}

func ComputeFastCompassBySecond(fastCompass metrics.FastCompass) (float64, float64, float64, float64, float64, float64) {
	bucketInterval := fastCompass.GetBucketInterval().Milliseconds() / 1000
	start := getNormalizedStartTime(bucketInterval)

	var totalCount int64 = 0
	var successCount int64 = 0
	var errorCount int64 = 0
	var totalRt int64 = 0
	var successRt int64 = 0
	var errorRt int64 = 0

	for k1, v1 := range fastCompass.GetMethodCountPerCategorySince(start) {
		if v2, ok := v1[start]; ok {
			totalCount += v2
			if k1 == "success" {
				successCount += v2
			} else {
				errorCount += v2
			}
		}
	}

	for k1, v1 := range fastCompass.GetMethodRtPerCategorySince(start) {
		if v2, ok := v1[start]; ok {
			totalRt += v2
			if k1 == "success" {
				successRt += v2
			} else {
				errorRt += v2
			}
		}
	}

	return rate(totalCount, bucketInterval), rate(successCount, bucketInterval), rate(errorCount, bucketInterval),
		rate(totalRt, totalCount), rate(successRt, successCount), rate(errorRt, errorCount)
}

func rate(data, interval int64) float64 {
	if interval == 0 {
		return 0
	}

	return float64(data) / float64(interval)
}

func getNormalizedStartTime(interval int64) int64 {
	return (time.Now().Unix() - interval) / interval * interval * 1000
}
