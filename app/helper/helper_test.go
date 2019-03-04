package helper

import (
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func TestNullTimeToString(t *testing.T) {
	a := mysql.NullTime{}

	if NullTimeToString(a, time.RFC3339) != "" {
		t.Fatal("Passing invalid NullTime should return empty string")
	}
	t.Log(a)
}

func TestStringToInt(t *testing.T) {
	a := "1"

	if StringToInt(a) != 1 {
		t.Fatal("Passing to string should converted to int")
	}
	t.Log(a)
}
