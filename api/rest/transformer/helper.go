package transformer

import "time"

func getAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	if now.YearDay() < dob.YearDay() {
		years--
	}

	return years
}
