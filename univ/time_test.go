package univ

import (
	"testing"
	"time"
)

func TestAgo(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"小于1分钟",
			args{time.Now().Add(-60 * time.Second)},
			"1 minute ago",
		},
		{
			"约2小时前",
			args{time.Now().Add(-120 * time.Minute)},
			"2 hour ago",
		},
		{
			"约1天前",
			args{time.Now().Add(-25 * time.Hour)},
			"1 day ago",
		},
		{
			"01-01 00:00",
			args{time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)},
			"01-01 00:00",
		},
		{
			"2016-02-02",
			args{time.Date(2016, 2, 2, 0, 0, 0, 0, time.Local)},
			"2016-02-02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeAgo(tt.args.t); got != tt.want {
				t.Errorf("Ago() = %v, want %v", got, tt.want)
			}
		})
	}
}
