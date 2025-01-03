package game

import (
	"testing"
	"time"
)

func TestPointValidation(t *testing.T) {
	testCases := []struct {
		name      string
		time      string
		expected  bool
		points    int
		isSpecial bool
	}{
		{
			name:      "Valid 1337 time",
			time:      "13:37:00",
			expected:  true,
			points:    60, // Maximum points at start of minute
			isSpecial: false,
		},
		{
			name:      "Invalid 1337 time",
			time:      "13:38:00",
			expected:  false,
			points:    0,
			isSpecial: false,
		},
		{
			name:      "Valid 420 time",
			time:      "16:20:00",
			expected:  true,
			points:    6, // 10% of 60 (max 1337 points)
			isSpecial: true,
		},
		{
			name:      "Invalid 420 time",
			time:      "16:21:00",
			expected:  false,
			points:    0,
			isSpecial: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeStr := tc.time
			parsedTime, _ := time.Parse("15:04:05", timeStr)

			var valid bool
			var points int

			if tc.isSpecial {
				valid = IsValidSpecialTime(parsedTime)
				points = CalculateSpecialPoints(parsedTime)
			} else {
				valid = IsValid1337Time(parsedTime)
				points = Calculate1337Points(parsedTime)
			}

			if valid != tc.expected {
				t.Errorf("Time validation failed for %s: got %v, want %v", timeStr, valid, tc.expected)
			}

			if points != tc.points {
				t.Errorf("Point calculation failed for %s: got %d, want %d", timeStr, points, tc.points)
			}
		})
	}
}

func TestEasterEggSystem(t *testing.T) {
	testCases := []struct {
		name           string
		datetime       time.Time
		expectedType   string
		expectedValid  bool
		expectedPoints int
	}{
		{
			name:           "Regular 420 time",
			datetime:       time.Date(2024, 1, 1, 16, 20, 0, 0, time.UTC),
			expectedType:   "420",
			expectedValid:  true,
			expectedPoints: 6, // 10% of 60 (max 1337 points)
		},
		{
			name:           "Special 420 on April 20th",
			datetime:       time.Date(2024, 4, 20, 16, 20, 0, 0, time.UTC),
			expectedType:   "420_special",
			expectedValid:  true,
			expectedPoints: 600, // 10x of 60 (max 1337 points)
		},
		{
			name:           "Regular 420 on different April day",
			datetime:       time.Date(2024, 4, 21, 16, 20, 0, 0, time.UTC),
			expectedType:   "420",
			expectedValid:  true,
			expectedPoints: 6,
		},
		{
			name:           "Invalid time",
			datetime:       time.Date(2024, 4, 20, 16, 21, 0, 0, time.UTC),
			expectedType:   "",
			expectedValid:  false,
			expectedPoints: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test type detection
			eggType := GetEasterEggType(tc.datetime)
			if eggType != tc.expectedType {
				t.Errorf("GetEasterEggType() = %v, want %v", eggType, tc.expectedType)
			}

			// Test validation
			valid := IsValidSpecialTime(tc.datetime)
			if valid != tc.expectedValid {
				t.Errorf("IsValidSpecialTime() = %v, want %v", valid, tc.expectedValid)
			}

			// Test point calculation
			points := CalculateSpecialPoints(tc.datetime)
			if points != tc.expectedPoints {
				t.Errorf("CalculateSpecialPoints() = %v, want %v", points, tc.expectedPoints)
			}
		})
	}
}

func Test1337PointDecay(t *testing.T) {
	testCases := []struct {
		name           string
		time           time.Time
		expectedValid  bool
		expectedPoints int
	}{
		{
			name:           "Exactly 13:37:00",
			time:           time.Date(2024, 1, 1, 13, 37, 0, 0, time.UTC),
			expectedValid:  true,
			expectedPoints: 60,
		},
		{
			name:           "13:37:30 (30 seconds in)",
			time:           time.Date(2024, 1, 1, 13, 37, 30, 0, time.UTC),
			expectedValid:  true,
			expectedPoints: 30,
		},
		{
			name:           "13:37:59 (last valid second)",
			time:           time.Date(2024, 1, 1, 13, 37, 59, 0, time.UTC),
			expectedValid:  true,
			expectedPoints: 1,
		},
		{
			name:           "13:37:10 (50 points remaining)",
			time:           time.Date(2024, 1, 1, 13, 37, 10, 0, time.UTC),
			expectedValid:  true,
			expectedPoints: 50,
		},
		{
			name:           "13:38:00 (invalid time)",
			time:           time.Date(2024, 1, 1, 13, 38, 0, 0, time.UTC),
			expectedValid:  false,
			expectedPoints: 0,
		},
		{
			name:           "13:36:59 (invalid time)",
			time:           time.Date(2024, 1, 1, 13, 36, 59, 0, time.UTC),
			expectedValid:  false,
			expectedPoints: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test validation
			valid := IsValid1337Time(tc.time)
			if valid != tc.expectedValid {
				t.Errorf("IsValid1337Time() = %v, want %v", valid, tc.expectedValid)
			}

			// Test point calculation
			points := Calculate1337Points(tc.time)
			if points != tc.expectedPoints {
				t.Errorf("Calculate1337Points() = %v, want %v", points, tc.expectedPoints)
			}
		})
	}
}
