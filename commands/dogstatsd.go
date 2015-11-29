package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"time"
)

func makeTags(exec, watchType, watchId, id string) []string {
	tags := make([]string, 4)
	execTag := fmt.Sprintf("exec:%s", exec)
	watchTypeTag := fmt.Sprintf("watchtype:%s", watchType)
	watchIdTag := fmt.Sprintf("watchid:%s", watchId)
	idTag := fmt.Sprintf("id:%s", id)
	tags = append(tags, execTag)
	tags = append(tags, watchTypeTag)
	tags = append(tags, watchIdTag)
	tags = append(tags, idTag)
	return tags
}

func StatsdRunTime(start time.Time, exec string, watchType string, watchId string, id string) {
	if DogStatsd {
		elapsed := time.Since(start)
		milliseconds := int64(elapsed / time.Millisecond)
		Log(fmt.Sprintf("dogstatsd='true' %s='%s' exec='%s' id='%s' elapsed='%s'", watchType, watchId, exec, id, elapsed), "debug")
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(exec, watchType, watchId, id)
		metricName := fmt.Sprintf("%s.time", MetricPrefix)
		statsd.Gauge(metricName, float64(milliseconds), tags)
	}
}

func StatsdDuplicate(watchType string, watchId string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := make([]string, 2)
		watchTypeTag := fmt.Sprintf("watchtype:%s", watchType)
		watchIdTag := fmt.Sprintf("watchid:%s", watchId)
		tags = append(tags, watchTypeTag)
		tags = append(tags, watchIdTag)
		metricName := fmt.Sprintf("%s.duplicate", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
}

func StatsdBlank(watchType string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := make([]string, 1)
		watchTypeTag := fmt.Sprintf("watchtype:%s", watchType)
		tags = append(tags, watchTypeTag)
		metricName := fmt.Sprintf("%s.blank", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
}
