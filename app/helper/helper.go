package helper

import (
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func NullTimeToString(i mysql.NullTime, format string) string {
	if i.Valid {
		return i.Time.Format(format)
	}
	return ""
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
