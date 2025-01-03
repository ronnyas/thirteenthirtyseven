package game

import (
	"testing"
	"time"
)

func Test_Calculate1337Points(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		want      int
	}{
		{
			name:      "exactly at 13:37:00",
			timestamp: time.Date(2024, 1, 1, 13, 37, 0, 0, time.UTC),
			want:      60,
		},
		{
			name:      "30 seconds in",
			timestamp: time.Date(2024, 1, 1, 13, 37, 30, 0, time.UTC),
			want:      30,
		},
		{
			name:      "last valid second",
			timestamp: time.Date(2024, 1, 1, 13, 37, 59, 0, time.UTC),
			want:      1,
		},
		{
			name:      "invalid time - wrong hour",
			timestamp: time.Date(2024, 1, 1, 14, 37, 0, 0, time.UTC),
			want:      0,
		},
		{
			name:      "invalid time - wrong minute",
			timestamp: time.Date(2024, 1, 1, 13, 38, 0, 0, time.UTC),
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Calculate1337Points(tt.timestamp); got != tt.want {
				t.Errorf("Calculate1337Points() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EasterEggPoints(t *testing.T) {
	tests := []struct {
		name       string
		timestamp  time.Time
		wantType   string
		wantValid  bool
		wantPoints int
	}{
		{
			name:       "regular 420 time",
			timestamp:  time.Date(2024, 1, 1, 16, 20, 0, 0, time.UTC),
			wantType:   "420",
			wantValid:  true,
			wantPoints: 6, // 10% of 60 points
		},
		{
			name:       "special 420 on April 20th",
			timestamp:  time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC),
			wantType:   "420_special",
			wantValid:  true,
			wantPoints: 600, // 10x regular points
		},
		{
			name:       "regular 420 on different April day",
			timestamp:  time.Date(2024, 4, 10, 16, 20, 0, 0, time.UTC),
			wantType:   "420",
			wantValid:  true,
			wantPoints: 6,
		},
		{
			name:       "invalid time",
			timestamp:  time.Date(2024, 1, 1, 16, 21, 0, 0, time.UTC),
			wantType:   "",
			wantValid:  false,
			wantPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := GetEasterEggType(tt.timestamp)
			if gotType != tt.wantType {
				t.Errorf("GetEasterEggType() = %v, want %v", gotType, tt.wantType)
			}

			gotValid := IsValidSpecialTime(tt.timestamp)
			if gotValid != tt.wantValid {
				t.Errorf("IsValidSpecialTime() = %v, want %v", gotValid, tt.wantValid)
			}

			gotPoints := CalculateSpecialPoints(tt.timestamp)
			if gotPoints != tt.wantPoints {
				t.Errorf("CalculateSpecialPoints() = %v, want %v", gotPoints, tt.wantPoints)
			}
		})
	}
}
