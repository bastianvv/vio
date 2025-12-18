package tmdb

import "strconv"

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func parseYear(date string) int {
	if len(date) < 4 {
		return 0
	}

	year, err := strconv.Atoi(date[:4])
	if err != nil {
		return 0
	}
	return year
}
