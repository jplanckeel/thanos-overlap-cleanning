package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

type metaData struct {
	Thanos  thanosMeta `json:"thanos"`
	Ulid    string     `json:"ulid"`
	MinTime int64      `json:"minTime"`
	MaxTime int64      `json:"maxTime"`
}

type thanosMeta struct {
	Labels     map[string]string      `json:"labels"`
	Downsample map[string]interface{} `json:"downsample"`
	Source     string                 `json:"source"`
}

func getThanosMeta(metaFile string) (*metaData, error) {
	var d *metaData
	if err := json.Unmarshal([]byte(metaFile), &d); err != nil {
		return nil, err
	}
	logrus.Debug(d.Thanos.Labels)
	return d, nil
}

func appendOverlap(res []string, meta *metaData, ulid string) []string {
	exist := false
	for _, r := range res {
		if r == meta.Ulid {
			exist = true
		}
	}
	if !exist {
		res = append(res, meta.Ulid)
	}
	return append(res, ulid)
}

// metaOverlap ...
func metaOverlap(allMetas map[string]*metaData, checkMeta *metaData) string {
	for ulid, meta := range allMetas {
		if ulid == checkMeta.Ulid {
			continue
		}
		if !reflect.DeepEqual(meta.Thanos.Labels, checkMeta.Thanos.Labels) {
			continue
		}
		if (checkMeta.MinTime >= meta.MinTime && checkMeta.MinTime <= meta.MaxTime) ||
			(checkMeta.MaxTime >= meta.MinTime && checkMeta.MaxTime <= meta.MaxTime) {
			return ulid
		}
	}
	return ""
}

func CheckOverlap(dryrun bool, accessKey string, secretKey string, bucketName string, region string, maxTime string, minTime string, labelsSelector string, cacheDir string) {
	c, _ := newClient("s3.fr-par.scw.cloud", bucketName, accessKey, secretKey, region, maxTime, minTime, labelsSelector)

	var meta *metaData
	allMetas := make(map[string]*metaData)
	metadataFiles := c.listMeta()
	for _, files := range metadataFiles {

		if checkCache(files, cacheDir) {
			o, _ := c.getObjectFileContent(files)
			// Write in ./data
			writeCache(o, files, cacheDir)
			meta, _ = getThanosMeta(o)
		} else {
			rc, err := readCache(files, cacheDir)
			if err != nil {
				logrus.Fatal(err)
			}
			meta, err = getThanosMeta(rc)
			if err != nil {
				logrus.Fatal(err)
			}
		}

		if meta != nil && filterMetaData(meta, maxTime, minTime, labelsSelector) {
			allMetas[meta.Ulid] = meta
		}

	}

	//res := make(map[string][]string)
	for object, checkMeta := range allMetas {
		logrus.WithField("object", object).Debug("listing object")
		if ulid := metaOverlap(allMetas, checkMeta); ulid != "" {
			fmt.Printf("labels: %v ulid: %v  maxtime: %v mintime: %v \n", checkMeta.Thanos.Labels, ulid, time.UnixMilli(checkMeta.MaxTime), time.UnixMilli(checkMeta.MinTime))
			//res[fmt.Sprintf("%v", checkMeta.Thanos.Labels)] = appendOverlap(res[fmt.Sprintf("%v", checkMeta.Thanos.Labels)], checkMeta, ulid)
		}
	}

}

// Check if metadata match with condition maxTime,minTime and labelsSelector
func filterMetaData(meta *metaData, maxTime string, minTime string, labelsSelector string) bool {
	var checkT bool = true
	var checkL bool = true
	labelsFilter := stringToMap(labelsSelector)
	logrus.Debug("labels selector: %s", labelsFilter)
	if labelsSelector != "" {
		if !reflect.DeepEqual(meta.Thanos.Labels, labelsFilter) {
			checkT = false
		}
	}
	if maxTime != "" || minTime != "" {
		if !checkTime(meta, minTime, maxTime) {
			checkT = false
		}
	}
	return checkT && checkL
}
