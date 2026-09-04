package platform

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type UptimeInformation struct {
	DurationSeconds uint64 `json:"durationSeconds"`
	BootTime        string `json:"bootTime,omitempty"`
}

func GetUptimeInformation(ctx context.Context) (*UptimeInformation, error) {
	info := &UptimeInformation{}
	var errs []error

	uptime, err := host.UptimeWithContext(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("duration: %w", err))
	} else {
		info.DurationSeconds = uptime
	}
	bootTime, err := host.BootTimeWithContext(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("boot time: %w", err))
	} else {
		info.BootTime = time.Unix(int64(bootTime), 0).Format(time.RFC3339)
	}

	return info, errors.Join(errs...)
}
