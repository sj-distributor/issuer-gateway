package schedule

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewRedisScheduler(t *testing.T) {

	tests := []struct {
		name string
		want IScheduler
	}{
		{
			name: "can created redis schedule",
			want: NewRedisScheduler([]string{"127.0.0.1:6379"}, "", "", "", 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.want; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisScheduler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisScheduler_StartAsync(t *testing.T) {
	type args struct {
		cron        string
		onExecuting func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "redis scheduler can running",
			args: struct {
				cron        string
				onExecuting func()
			}{cron: "*/1 * * * *", onExecuting: func() {
				fmt.Println("do action")
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedisScheduler([]string{"127.0.0.1:6379"}, "", "", "", 0)
			if err := r.StartAsync(tt.args.cron, tt.args.onExecuting); (err != nil) != tt.wantErr {
				t.Errorf("StartAsync() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(time.Duration(1) * time.Minute)
		})
	}
}

func TestRedisScheduler_Stop(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "can stop redis schedule",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedisScheduler([]string{"127.0.0.1:6379"}, "", "", "", 0)

			err := r.StartAsync("*/1 * * * *", func() {
				fmt.Println("on action...")
			})
			if err != nil {
				t.Errorf("redis schedule starting failed..")
			}
			if err := r.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(time.Duration(2) * time.Second)
		})
	}
}
