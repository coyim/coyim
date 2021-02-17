package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
)

type timeTranslator struct {
	dayIndex, monthIndex, yearIndex int
	applies                         func(time.Time, timeTranslator) bool
	formatter                       func(time.Time) string
}

var timeTranslators []timeTranslator

func getTimeTranslators() []timeTranslator {
	if len(timeTranslators) == 0 {
		timeTranslators = []timeTranslator{
			timeTranslator{
				0, 0, 0,
				checkTimeToday,
				func(t time.Time) string {
					return i18n.Local("Today")
				},
			},
			timeTranslator{
				-1, 0, 0,
				checkTimeAfter,
				func(t time.Time) string {
					return i18n.Local("Yesterday")
				},
			},
			timeTranslator{
				-2, 0, 0,
				checkTimeAfter,
				func(t time.Time) string {
					return i18n.Local("Two days ago")
				},
			},
			timeTranslator{
				-3, 0, 0,
				checkTimeAfter,
				func(t time.Time) string {
					return i18n.Local("Three days ago")
				},
			},
			timeTranslator{
				-4, 0, 0,
				checkTimeAfter,
				func(t time.Time) string {
					return i18n.Local("Four days ago")
				},
			},
			timeTranslator{
				0, -1, 0,
				checkTimeAfter,
				func(t time.Time) string {
					return timeToFriendlyDate(t)
				},
			},
		}
	}

	return timeTranslators
}

func timeToFriendlyString(t time.Time) string {
	for _, tt := range getTimeTranslators() {
		if tt.applies(t, tt) {
			return tt.formatter(t)
		}
	}

	return t.Format(time.ANSIC)
}

func timeToFriendlyDate(t time.Time) string {
	return i18n.Localf("%s, %s %d, %d", localizedWeekday(t.Weekday()), localizedMonth(t.Month()), t.Day(), t.Year())
}

func localizedWeekday(wd time.Weekday) string {
	switch wd {
	case time.Monday:
		return i18n.Local("Monday")
	case time.Thursday:
		return i18n.Local("Thursday")
	case time.Wednesday:
		return i18n.Local("Wednesday")
	case time.Tuesday:
		return i18n.Local("Tuesday")
	case time.Friday:
		return i18n.Local("Friday")
	case time.Saturday:
		return i18n.Local("Saturday")
	case time.Sunday:
		return i18n.Local("Sunday")
	}

	return ""
}

func localizedMonth(m time.Month) string {
	switch m {
	case time.January:
		return i18n.Local("January")
	case time.February:
		return i18n.Local("February")
	case time.March:
		return i18n.Local("March")
	case time.April:
		return i18n.Local("April")
	case time.May:
		return i18n.Local("May")
	case time.June:
		return i18n.Local("June")
	case time.July:
		return i18n.Local("July")
	case time.August:
		return i18n.Local("August")
	case time.September:
		return i18n.Local("September")
	case time.October:
		return i18n.Local("October")
	case time.November:
		return i18n.Local("November")
	case time.December:
		return i18n.Local("December")
	}

	return ""
}

func getTodayFromStart() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
}

func checkTimeToday(t time.Time, tt timeTranslator) bool {
	return areTheSameDateInUTC(t, time.Now())
}

func checkTimeAfter(t time.Time, tt timeTranslator) bool {
	return t.After(getTodayFromStart().AddDate(tt.yearIndex, tt.monthIndex, tt.dayIndex))
}

func areTheSameDateInUTC(d1, d2 time.Time) bool {
	t1y, t1m, t1d := d1.In(time.UTC).Date()
	t2y, t2m, t2d := d2.In(time.UTC).Date()

	return t1d == t2d && t1m == t2m && t1y == t2y
}

func elapsedFriendlyTime(t time.Time) string {
	diff := time.Now().Sub(t)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes())
	seconds := int(diff.Seconds())

	switch {
	case hours > 24:
		return timeToFriendlyString(t)
	case hours == 1:
		return i18n.Local("An hour ago")
	case hours > 1 && hours <= 10:
		return i18n.Local("Few hours ago")
	case hours > 10:
		return i18n.Localf("%v hours ago", hours)
	case minutes == 1:
		return i18n.Local("A minute ago")
	case minutes > 1 && minutes <= 10:
		return i18n.Local("A few minutes ago")
	case minutes > 10:
		return i18n.Localf("%v minutes ago", minutes)
	case seconds == 1:
		return i18n.Localf("A second ago")
	case seconds > 1 && seconds <= 10:
		return i18n.Local("A few seconds ago")
	case seconds > 10:
		return i18n.Localf("%v seconds ago", seconds)
	default:
		return i18n.Local("Now")
	}
}
