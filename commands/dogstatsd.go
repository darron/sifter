package commands

import (
	"fmt"
	"github.com/PagerDuty/godspeed"
	"time"
)

func makeTags(event, exec, ltime string) []string {
	tags := make([]string, 3)
	eventTag := fmt.Sprintf("event:%s", event)
	execTag := fmt.Sprintf("exec:%s", exec)
	lTimeTag := fmt.Sprintf("ltime:%s", ltime)
	tags = append(tags, eventTag)
	tags = append(tags, execTag)
	tags = append(tags, lTimeTag)
	return tags
}

func StatsdRunTime(start time.Time, event string, exec string, ltime int64) {
	elapsed := time.Since(start)
	milliseconds := int64(elapsed / time.Millisecond)
	Log(fmt.Sprintf("dogstatsd='true' event='%s' exec='%s' ltime='%d' elapsed='%s'", event, exec, ltime, elapsed), "info")
	statsd, _ := godspeed.NewDefault()
	defer statsd.Conn.Close()
	tags := makeTags(event, exec, string(ltime))
	statsd.Gauge("sifter.time", float64(milliseconds), tags)
}
