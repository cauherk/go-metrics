package main

import (
	"fmt"
	"github.com/cauherk/go-metrics"
	"github.com/cauherk/go-metrics/impl"
	"time"
)

func main() {
	mr := impl.NewMetricRegistry(1000)
	mn := metrics.NewMetricName("hsf", nil, metrics.Major)
	fc := mr.GetFastCompass(mn)
	fmt.Println(fc)

	for i := 0; i < 8; i++ {
		go func() {
			for j := 0; j < 100000; j++ {
				fc.Record(22*time.Millisecond, "success")
				time.Sleep(200 * time.Millisecond)
			}
		}()
	}

	go func() {
		for i := 0; i < 100000; i++ {
			time.Sleep(1 * time.Second)

			bucketInterval := fc.GetBucketInterval().Milliseconds() / 1000
			start := getNormalizedStartTime(bucketInterval)

			var totalCount int64 = 0
			var successCount int64 = 0
			var errorCount int64 = 0
			var totalRt int64 = 0
			var successRt int64 = 0
			var errorRt int64 = 0

			for k1, v1 := range fc.GetMethodCountPerCategorySince(start) {
				if v2, ok := v1[start]; ok {
					totalCount += v2
					if k1 == "success" {
						successCount += v2
					} else {
						errorCount += v2
					}
				}
			}

			for k1, v1 := range fc.GetMethodRtPerCategorySince(start) {
				if v2, ok := v1[start]; ok {
					totalRt += v2
					if k1 == "success" {
						successRt += v2
					} else {
						errorRt += v2
					}
				}
			}

			fmt.Printf("key: %s , qps: %.2f , rt: %.2f \n", mn.Key, rate(totalCount, bucketInterval), rate(totalRt, totalCount))
		}
	}()
	select {}
}

func rate(data, interval int64) float64 {
	if interval == 0 {
		return 0
	}

	return float64(data) / float64(interval)
}

/**
 *
 * @param current
 * @param interval
 * @return
 */
func getNormalizedStartTime(interval int64) int64 {
	return (time.Now().Unix() - interval) / interval * interval * 1000
}
