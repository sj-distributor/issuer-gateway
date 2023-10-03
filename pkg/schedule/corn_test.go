package schedule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCronParse(t *testing.T) {
	type args struct {
		corn        string
		loc         *time.Location
		withSeconds bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "can create Cron Parser", args: struct {
			corn        string
			loc         *time.Location
			withSeconds bool
		}{corn: "*/1 * * * *", loc: time.UTC, withSeconds: false},
			wantErr: false,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CronParse(tt.args.corn, tt.args.loc, tt.args.withSeconds)

			if (err != nil) != tt.wantErr {
				t.Errorf("CronParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			utc := time.Now().UTC()
			next := got.Next(utc)

			utc.Add(time.Duration(1) * time.Minute).Unix()

			assert.True(t, next.Unix() > utc.Unix())
		})
	}
}
