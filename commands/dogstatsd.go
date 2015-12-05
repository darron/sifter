package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"time"
)

// StatsdRunTime sends timing metrics via Dogstatsd.
func StatsdRunTime(start time.Time, exec string, watchType string, watchID string, id string) {
	elapsed := time.Since(start)
	if DogStatsd {
		milliseconds := int64(elapsed / time.Millisecond)
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchID, exec, id)
		metricName := fmt.Sprintf("%s.time", MetricPrefix)
		statsd.Gauge(metricName, float64(milliseconds), tags)
	}
	Log(fmt.Sprintf("dogstatsd='%t' %s='%s' exec='%s' id='%s' elapsed='%s'", DogStatsd, watchType, watchID, exec, id, elapsed), "debug")
}

// StatsdDuplicate sends Dogstatsd metrics whenever there's a duplicate watch.
func StatsdDuplicate(watchType string, watchID string) {
	if DogStatsd {
		statsd, _ := godspeed.NewDefault()
		defer statsd.Conn.Close()
		tags := makeTags(watchType, watchID, "", "")
		metricName := fmt.Sprintf("%s.duplicate", MetricPrefix)
		statsd.Incr(metricName, tags)
	}
	Log(fmt.Sprintf("dogstatsd='%t' %s='%s' action='duplicate'", DogStatsd, watchType, watchID), "debug")
}

// StatsdBlank sends Dogstatsd metrics whenever there's a blank watch.
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

func makeTags(watchType, watchID, exec, id string) []string {
	var tags []string
	if watchType != "" {
		watchTypeTag := fmt.Sprintf("watchtype:%s", watchType)
		tags = append(tags, watchTypeTag)
	}
	if watchID != "" {
		watchIDTag := fmt.Sprintf("watchid:%s", watchID)
		tags = append(tags, watchIDTag)
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
