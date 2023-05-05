package pkg

import (
	"fmt"
	"strings"
	"time"
)

// Convert string to map[string]string
func stringToMap(s string) map[string]string {
	entries := strings.Split(s, ",")

	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, ":")
		m[parts[0]] = parts[1]
	}
	return m
}

func checkTime(meta *metaData, minTime string, maxTime string) bool {
	layout := "2006-01-02T15:04:05Z"
	minT, _ := time.Parse(layout, minTime)
	maxT, _ := time.Parse(layout, maxTime)
	metaMinT := time.UnixMilli(meta.MinTime)
	metaMaxT := time.UnixMilli(meta.MaxTime)

	if minT.Before(metaMinT) && maxT.After(metaMaxT) {
		fmt.Println("meta is true in time")
		return true
	}
	//fmt.Printf("meta is false in time minT:%v MaxT:%v metaMinT:%v metaMaxT%v", minT, maxT, metaMinT, metaMaxT)
	return false
}
