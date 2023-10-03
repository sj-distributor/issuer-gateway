package schedule

import (
	"github.com/go-co-op/gocron"
	"time"
)

type CronScheduler struct {
	s   *gocron.Scheduler
	job *gocron.Job
}

func NewCronScheduler() IScheduler {
	return &CronScheduler{
		s: gocron.NewScheduler(time.UTC),
	}
}

func (c *CronScheduler) StartAsync(cron string, onExecuting func()) error {

	job, err := c.s.Cron(cron).Do(onExecuting)

	if err != nil {
		return err
	}

	c.job = job

	c.s.StartAsync()
	return nil
}

func (c *CronScheduler) Stop() error {
	c.s.Stop()
	return nil
}
