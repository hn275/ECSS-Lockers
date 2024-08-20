package time

import (
	stdtime "time"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	loc *stdtime.Location
)

func init() {
	var err error
	loc, err = stdtime.LoadLocation("Canada/Pacific")
	if err != nil {
		logger.Fatal(err)
	}
}

func NextExpiryDate(now stdtime.Time) stdtime.Time {
	const (
		expHour    int = 23
		expMin     int = 59
		expSecond  int = 59
		expNSecond int = 0
	)

	currentMonth := now.Month()

	var (
		expMonth stdtime.Month
		expDate  int
		expYear  int = now.Year()
	)

	if currentMonth <= 4 { // spring term
		expMonth = stdtime.Month(4)
		expDate = 30

	} else if currentMonth <= 8 { // summer term
		expMonth = stdtime.Month(8)
		expDate = 31

	} else { // winter term
		expMonth = stdtime.Month(12)
		expDate = 31
	}

	return stdtime.Date(
		expYear, expMonth, expDate,
		expHour, expMin, expSecond, expNSecond, loc)
}

func Now() stdtime.Time {
	return stdtime.Now().In(loc)
}
