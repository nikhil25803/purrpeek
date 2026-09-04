package platform

import (
	"fmt"
	"time"
)

type TimeInfo struct {
	CurrentTime string `json:"currentTime"`
	TimeZone    string `json:"timeZone"`
	UTCOffset   string `json:"utcOffset"`
}

func GetTimeInformation() *TimeInfo {
	return timeInformation(time.Now())
}

func timeInformation(now time.Time) *TimeInfo {
	zone, offset := now.Zone()
	sign := '+'
	if offset < 0 {
		sign = '-'
		offset = -offset
	}
	return &TimeInfo{
		CurrentTime: now.Format(time.RFC3339),
		TimeZone:    zone,
		UTCOffset:   fmt.Sprintf("%c%02d:%02d", sign, offset/3600, offset%3600/60),
	}
}
