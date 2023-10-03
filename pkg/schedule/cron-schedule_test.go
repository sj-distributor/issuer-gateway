package schedule

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/assert/v2"
	"reflect"
	"testing"
	"time"
)

func TestCronSchedule_StartAsync(t *testing.T) {
	type fields struct {
		s   *gocron.Scheduler
		job *gocron.Job
	}
	type args struct {
		cron string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "can test schedule is running",
			fields: struct {
				s   *gocron.Scheduler
				job *gocron.Job
			}{s: gocron.NewScheduler(time.UTC), job: nil},
			args: struct{ cron string }{cron: "*/1 * * * *"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CronScheduler{
				s:   tt.fields.s,
				job: tt.fields.job,
			}
			if err := c.StartAsync(tt.args.cron, func() {
				fmt.Println("on Executing...")
			}); (err != nil) != tt.wantErr {
				t.Errorf("StartAsync() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotEqual(t, c.job, nil)
		})
	}
}

func TestCronSchedule_Stop(t *testing.T) {
	type fields struct {
		s   *gocron.Scheduler
		job *gocron.Job
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "can stop schedule",
			fields: struct {
				s   *gocron.Scheduler
				job *gocron.Job
			}{s: gocron.NewScheduler(time.UTC), job: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CronScheduler{
				s:   tt.fields.s,
				job: tt.fields.job,
			}
			err := c.StartAsync("*/1 * * * *", func() {
				fmt.Println("on Executing...")
			})
			if err != nil {
				t.Errorf("start is error %s", err)
			}
			time.Sleep(time.Second * 2)
			if err := c.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCronSchedule(t *testing.T) {
	tests := []struct {
		name string
		want IScheduler
	}{
		{name: "can created schedule", want: NewCronScheduler()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.want; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCronSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}
