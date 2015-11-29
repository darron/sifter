package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"time"
)

func StatsdRunTime(start time.Time, exec string, watchType string, watchId string, id string) {
	if DogStatsd {
		elapsed := time.Since(start)
		milliseconds := int64(elapsed / time.Millisecond)
		Log(fmt.Sprintf("dogstatsd='true' %s='%s' exec='%s' id='%s' elapsed='%s'", watchType, watchId, exec, id, elapsed), "debug")
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchId, exec, id)
		metricName := fmt.Sprintf("%s.time", MetricPrefix)
		statsd.Gauge(metricName, float64(milliseconds), tags)
	}
}

func StatsdDuplicate(watchType string, watchId string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchId, "", "")
		metricName := fmt.Sprintf("%s.duplicate", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
}

func StatsdBlank(watchType string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, "", "", "")
		metricName := fmt.Sprintf("%s.blank", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
}

func makeTags(watchType, watchId, exec, id string) []string {
	tags := make([]string, 0)
	if watchType != "" {
		watchTypeTag := fmt.Sprintf("watchtype:%s", watchType)
		tags = append(tags, watchTypeTag)
	}
	if watchId != "" {
		watchIdTag := fmt.Sprintf("watchid:%s", watchId)
		tags = append(tags, watchIdTag)
	}
	if exec != "" {
		execTag := fmt.Sprintf("exec:%s", exec)
		tags = append(tags, execTag)
	}
	if id != "" {
		idTag := fmt.Sprintf("id:%s", id)
		tags = append(tags, idTag)
	}
	return tags
}
