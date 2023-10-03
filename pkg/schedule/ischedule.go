package schedule

type IScheduler interface {
	StartAsync(cron string, onExecuting func()) error
	Stop() error
}
