package dateutil

import "time"

func Today() (time.Time, time.Time) {

	today := time.Now()

	start := time.Date(
		today.Year(),
		today.Month(),
		today.Day(),
		0,
		0,
		0,
		0,
		today.Location(),
	)

	end := start.AddDate(0, 0, 1)

	return start, end

}

func Yesterday() (time.Time, time.Time) {

	start, end := Today()

	return start.AddDate(0, 0, -1), end.AddDate(0, 0, -1)

}

func CurrentMonth() (time.Time, time.Time) {

	now := time.Now()

	start := time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		now.Location(),
	)

	end := time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		now.Location(),
	)

	return start, end

}

func PreviousMonth() (time.Time, time.Time) {

	start, end := CurrentMonth()

	return start.AddDate(0, -1, 0), end.AddDate(0, -1, 0)

}
