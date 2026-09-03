package platform

import (
	"time"

	utils "github.com/nikhil25803/purrpeek/internal/utils"
	"github.com/shirou/gopsutil/v3/host"
)

type UptimeInformation struct {
	UptimeDuration string
	BootTime       string
}

func GetUptimeInformation() (*UptimeInformation, error) {

	uptimeSeconds, err := host.Uptime()
	if err != nil {
		return nil, err
	}

	bootTime, err := host.BootTime()
	if err != nil {
		return nil, err
	}

	uptimeDuration := time.Duration(uptimeSeconds) * time.Second

	bootDuration := time.Unix(int64(bootTime), 0)

	return &UptimeInformation{
		UptimeDuration: utils.FormatDuration(uptimeDuration),
		BootTime:       bootDuration.Format("2006-01-02 15:04:05"),
	}, nil
}
