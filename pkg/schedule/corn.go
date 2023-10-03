package schedule

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func CronParse(corn string, loc *time.Location, withSeconds bool) (cron.Schedule, error) {

	withLocation := fmt.Sprintf("CRON_TZ=%s %s", loc.String(), corn)

	var (
		cronSchedule cron.Schedule
		err          error
	)

	if withSeconds {
		p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		cronSchedule, err = p.Parse(withLocation)
	} else {
		cronSchedule, err = cron.ParseStandard(withLocation)
	}

	return cronSchedule, err
}
