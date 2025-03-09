package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func FormatTime(pgTime pgtype.Time) string {
	if !pgTime.Valid {
		return "" // Or some other default value
	}

	microseconds := pgTime.Microseconds
	hours := microseconds / (3600 * 1e6)
	microseconds %= 3600 * 1e6
	minutes := microseconds / (60 * 1e6)
	microseconds %= 60 * 1e6
	seconds := microseconds / 1e6

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
