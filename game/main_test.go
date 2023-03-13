package game

import (
	"testing"
	"time"
)

func Test_calculatePointsFromTimestamp(t *testing.T) {
	type args struct {
		timestamp time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				timestamp: time.Now(),
			},
			want: 60 - time.Now().Second(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePointsFromTimestamp(tt.args.timestamp); got != tt.want {
				t.Errorf("calculatePointsFromTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}