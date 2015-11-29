package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"time"
)

func StatsdRunTime(start time.Time, exec string, watchType string, watchId string, id string) {
	elapsed := time.Since(start)
	if DogStatsd {
		milliseconds := int64(elapsed / time.Millisecond)
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchId, exec, id)
		metricName := fmt.Sprintf("%s.time", MetricPrefix)
		statsd.Gauge(metricName, float64(milliseconds), tags)
	}
	Log(fmt.Sprintf("dogstatsd='%t' %s='%s' exec='%s' id='%s' elapsed='%s'", DogStatsd, watchType, watchId, exec, id, elapsed), "debug")
}

func StatsdDuplicate(watchType string, watchId string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchId, "", "")
		metricName := fmt.Sprintf("%s.duplicate", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
	Log(fmt.Sprintf("dogstatsd='%t' %s='%s' action='duplicate'", DogStatsd, watchType, watchId), "debug")
}

func StatsdBlank(watchType string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, "", "", "")
		metricName := fmt.Sprintf("%s.blank", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
	Log(fmt.Sprintf("dogstatsd='%t' watchType='%s' action='blank'", DogStatsd, watchType), "debug")
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
