package leet

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/ronnyas/thirteenthirtyseven/language"
	"golang.org/x/exp/slices"
)

type Streak struct {
	UserID    string
	StartTime string
	EndTime   string
}
type Point struct {
	UserID    string
	Timestamp string
	Points    int
}

func (s *Streak) Duration() int {
	layout := "2006-01-02 15:04:05.999999999-07:00"

	start, err := time.Parse(layout, s.StartTime)
	if err != nil {
		panic(err)
	}

	end, err := time.Parse(layout, s.EndTime)
	if err != nil {
		panic(err)
	}

	return int(end.Sub(start).Round(24*time.Hour).Hours()/24) + 1
}

// backfill streaks from points, since we didn't have streaks before.
// this is a one-time thing. keeping it for reference.
func BackfillStreaks(db *sql.DB) error {
	rows, err := db.Query(`
		select user_id, timestamp from points
		order by user_id, timestamp
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var streak Streak
	var streaks []Streak

	for rows.Next() {
		var userID string
		var timestamp string

		if err := rows.Scan(&userID, &timestamp); err != nil {
			return err
		}

		if streak.UserID == "" {
			log.Printf(language.GetTranslation("streak_new"), userID)
			streak.UserID = userID
			streak.StartTime = timestamp
			streak.EndTime = timestamp

			streaks = append(streaks, streak)
		} else if streak.UserID == userID {
			// check if streak.EndTime + 1 day == timestamp
			layout := "2006-01-02 15:04:05.999999999-07:00"
			end, err := time.Parse(layout, strings.TrimSpace(streak.EndTime))
			if err != nil {
				panic(err)
			}
			// remove time from end
			end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

			t, err := time.Parse(layout, strings.TrimSpace(timestamp))
			if err != nil {
				panic(err)
			}

			// remove time from t
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

			if end.AddDate(0, 0, 1).Equal(t) {
				log.Printf(language.GetTranslation("streak_continue"), userID)
				streak.EndTime = timestamp
			} else {
				log.Printf(language.GetTranslation("streak_new_because"), userID, end.AddDate(0, 0, 1), t)
				streaks = append(streaks, streak)
				streak.UserID = userID
				streak.StartTime = timestamp
				streak.EndTime = timestamp
			}
		} else {
			log.Printf(language.GetTranslation("streak_new"), userID)
			streak.UserID = userID
			streak.StartTime = timestamp
			streak.EndTime = timestamp
			streaks = append(streaks, streak)
		}
	}

	if streak.UserID != "" {
		streaks = append(streaks, streak)
	}

	for _, streak := range streaks {
		// check if streak is 3 days or more
		if streak.Duration() < Config.StreakDays {
			continue
		}

		_, err := db.Exec(`
			insert into streaks (user_id, start_time, end_time)
			values (?, ?, ?)
		`, streak.UserID, streak.StartTime, streak.EndTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// get streaks that are currently active
func GetActiveStreaks(db *sql.DB) ([]Streak, error) {
	// get all streaks that have end time == today
	timeNow := time.Now()
	var timeString string
	if timeNow.Hour() < 13 || (timeNow.Hour() == 13 && timeNow.Minute() < 38) {
		timeString = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	} else {
		timeString = time.Now().Format("2006-01-02")
	}

	rows, err := db.Query(`
		select user_id, start_time, end_time from streaks
		where end_time like ?
	`, timeString+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streaks []Streak

	for rows.Next() {
		var streak Streak

		if err := rows.Scan(&streak.UserID, &streak.StartTime, &streak.EndTime); err != nil {
			return nil, err
		}

		streaks = append(streaks, streak)
	}

	return streaks, nil
}

func GetTodaysPoints(db *sql.DB) ([]Point, error) {
	rows, err := db.Query(`
		select user_id, timestamp, points from points
		where timestamp like ?
	`, time.Now().Format("2006-01-02")+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []Point

	for rows.Next() {
		var point Point

		if err := rows.Scan(&point.UserID, &point.Timestamp, &point.Points); err != nil {
			return nil, err
		}

		points = append(points, point)
	}

	return points, nil
}

func RefreshActiveStreak(db *sql.DB, streak Streak, time string) error {
	sqlStatement := `
		update streaks
		set end_time = ?
		where user_id = ? and start_time = ? and end_time = ?
	`

	_, err := db.Exec(sqlStatement, time, streak.UserID, streak.StartTime, streak.EndTime)
	if err != nil {
		return err
	}
	return nil
}

func CreateStreak(db *sql.DB, streak Streak) error {
	sqlStatement := `
		insert into streaks (user_id, start_time, end_time)
		values (?, ?, ?)
	`
	_, err := db.Exec(sqlStatement, streak.UserID, streak.StartTime, streak.EndTime)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAllStreaks(db *sql.DB) (new []Streak, broken []Streak, error error) {
	// get all streaks that have end time == yesterday
	rows, err := db.Query(`
		select user_id, start_time, end_time from streaks
		where end_time like ?
	`, time.Now().AddDate(0, 0, -1).Format("2006-01-02")+"%")

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	var streaks []Streak

	for rows.Next() {
		var streak Streak

		if err := rows.Scan(&streak.UserID, &streak.StartTime, &streak.EndTime); err != nil {
			return nil, nil, err
		}

		streaks = append(streaks, streak)
	}

	// get all points that have been given today
	points, err := GetTodaysPoints(db)
	if err != nil {
		return nil, nil, err
	}

	var foundUsers []string
	var newStreaks []Streak

	for _, point := range points {
		for _, streak := range streaks {
			if streak.UserID == point.UserID {
				err := RefreshActiveStreak(db, streak, point.Timestamp)
				if err != nil {
					return nil, nil, err
				}

				foundUsers = append(foundUsers, point.UserID)

				break
			}
		}

		// if user doesn't have a streak, create one if they have received points the lasst three days
		if !slices.Contains(foundUsers, point.UserID) {
			rows2, err := db.Query(`
				select user_id, timestamp from points
				where user_id = ? and timestamp > ?
			`, point.UserID, time.Now().AddDate(0, 0, 0-Config.StreakDays).Format("2006-01-02")+"%")

			if err != nil {
				return nil, nil, err
			}
			defer rows2.Close()

			var points []Point

			for rows2.Next() {
				var point Point

				if err := rows2.Scan(&point.UserID, &point.Timestamp); err != nil {
					return nil, nil, err
				}

				points = append(points, point)
			}

			if len(points) >= 3 {
				// create streak
				streak := Streak{
					UserID:    point.UserID,
					StartTime: points[0].Timestamp,
					EndTime:   point.Timestamp,
				}

				err := CreateStreak(db, streak)
				if err != nil {
					return nil, nil, err
				}

				newStreaks = append(newStreaks, streak)
			}
		}
	}

	// check broken streaks
	var brokenStreaks []Streak
	for _, streak := range streaks {
		if !slices.Contains(foundUsers, streak.UserID) {
			brokenStreaks = append(brokenStreaks, streak)
		}
	}

	return newStreaks, brokenStreaks, nil
}
