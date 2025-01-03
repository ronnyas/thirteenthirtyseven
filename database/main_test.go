package database

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestConnect(t *testing.T) {
	type args struct {
		databasePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				databasePath: "test.db",
			},
			want:    &sql.DB{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.args.databasePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetupDatabaseSchema(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := SetupDatabaseSchema(db); err != nil {
		t.Errorf("SetupDatabaseSchema() error = %v", err)
	}

	// Verify tables exist
	tables := []string{"points", "special_points", "streaks"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestAddPoints(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := SetupDatabaseSchema(db); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	// Test adding points
	testTime := time.Date(2024, 1, 1, 13, 37, 0, 0, time.UTC)
	err = AddPoints(db, "user1", testTime, 60)
	if err != nil {
		t.Errorf("AddPoints() error = %v", err)
	}

	// Verify points were added
	var points int
	err = db.QueryRow("SELECT points FROM points WHERE user_id = ?", "user1").Scan(&points)
	if err != nil {
		t.Errorf("Failed to query points: %v", err)
	}
	if points != 60 {
		t.Errorf("Expected 60 points, got %d", points)
	}
}

func TestAddSpecialPoints(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := SetupDatabaseSchema(db); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	// Test adding special points
	testTime := time.Date(2024, 1, 1, 16, 20, 0, 0, time.UTC)
	err = AddSpecialPoints(db, "user1", testTime, 6, "420")
	if err != nil {
		t.Errorf("AddSpecialPoints() error = %v", err)
	}

	// Verify points were added
	var points int
	var eggType string
	err = db.QueryRow("SELECT points, type FROM special_points WHERE user_id = ?", "user1").Scan(&points, &eggType)
	if err != nil {
		t.Errorf("Failed to query special points: %v", err)
	}
	if points != 6 {
		t.Errorf("Expected 6 points, got %d", points)
	}
	if eggType != "420" {
		t.Errorf("Expected type '420', got %s", eggType)
	}
}

func TestGetTotalPoints(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := SetupDatabaseSchema(db); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO points (timestamp, user_id, points) VALUES 
		('2024-01-01 13:37:00', 'user1', 60),
		('2024-01-02 13:37:00', 'user1', 60);
		
		INSERT INTO special_points (timestamp, user_id, points, type) VALUES 
		('2024-01-01 16:20:00', 'user1', 6, '420'),
		('2024-01-02 16:20:00', 'user1', 6, '420');
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test total points calculation
	total, err := GetTotalPoints(db, "user1")
	if err != nil {
		t.Errorf("GetTotalPoints() error = %v", err)
	}

	expectedTotal := 132 // 2 * 60 (1337 points) + 2 * 6 (420 points)
	if total != expectedTotal {
		t.Errorf("GetTotalPoints() = %v, want %v", total, expectedTotal)
	}

	// Test non-existent user
	total, err = GetTotalPoints(db, "nonexistent")
	if err != nil {
		t.Errorf("GetTotalPoints() error for nonexistent user = %v", err)
	}
	if total != 0 {
		t.Errorf("GetTotalPoints() for nonexistent user = %v, want 0", total)
	}
}

func TestGetPointsForDay(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := SetupDatabaseSchema(db); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}

	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO points (timestamp, user_id, points) VALUES 
		('2024-01-01 13:37:00', 'user1', 60);
		
		INSERT INTO special_points (timestamp, user_id, points, type) VALUES 
		('2024-01-01 16:20:00', 'user1', 6, '420');
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test points for specific day
	total, err := GetPointsForDay(db, "user1", testDate)
	if err != nil {
		t.Errorf("GetPointsForDay() error = %v", err)
	}

	expectedTotal := 66 // 60 (1337 points) + 6 (420 points)
	if total != expectedTotal {
		t.Errorf("GetPointsForDay() = %v, want %v", total, expectedTotal)
	}

	// Test different day (should be 0)
	differentDate := testDate.AddDate(0, 0, 1)
	total, err = GetPointsForDay(db, "user1", differentDate)
	if err != nil {
		t.Errorf("GetPointsForDay() error = %v", err)
	}
	if total != 0 {
		t.Errorf("GetPointsForDay() for different day = %v, want 0", total)
	}
}
