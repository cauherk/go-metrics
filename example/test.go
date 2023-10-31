package main

import (
	"fmt"
	"github.com/cauherk/go-metrics/util"
	"time"
)

func main() {
	mr := util.NewOrGetMetricRegistry(1000)
	mn := util.NewMetricName("hsf", nil)
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

			tc, _, _, tt, _, _ := util.ComputeFastCompassBySecond(fc)
			fmt.Printf("key: %s , qps: %.2f , rt: %.2f \n", mn.Key, tc, tt)
		}
	}()
	select {}
}
