package common

import "time"

// IsUSHoliday checks whether a given date is a major public holiday in the United States.
// These are key holidays that can impact financial markets, including Forex.
// Includes both fixed-date and floating holidays.
//
// List of holidays covered:
// - New Year's Day: January 1
// - Martin Luther King Jr. Day: 3rd Monday of January
// - Presidents' Day (Washington's Birthday): 3rd Monday of February
// - Memorial Day: Last Monday of May
// - Independence Day: July 4
// - Labor Day: 1st Monday of September
// - Columbus Day: 2nd Monday of October
// - Veterans Day: November 11
// - Thanksgiving Day: 4th Thursday of November
// - Christmas Day: December 25
func IsUSHoliday(date time.Time) bool {
	month := date.Month()
	day := date.Day()

	// Fixed-date holidays
	switch {
	case month == time.January && day == 1:
		return true // New Year's Day
	case month == time.July && day == 4:
		return true // Independence Day
	case month == time.November && day == 11:
		return true // Veterans Day
	case month == time.December && day == 25:
		return true // Christmas Day
	}

	// Floating holidays
	switch {
	case isNthWeekdayOfMonth(date, time.Monday, 3): // MLK Day: 3rd Monday of January
		if month == time.January {
			return true
		}
	case isNthWeekdayOfMonth(date, time.Monday, 3): // Presidents' Day: 3rd Monday of February
		if month == time.February {
			return true
		}
	case isLastWeekdayOfMonth(date, time.Monday): // Memorial Day: last Monday of May
		if month == time.May {
			return true
		}
	case isNthWeekdayOfMonth(date, time.Monday, 1): // Labor Day: 1st Monday of September
		if month == time.September {
			return true
		}
	case isNthWeekdayOfMonth(date, time.Monday, 2): // Columbus Day: 2nd Monday of October
		if month == time.October {
			return true
		}
	case isNthWeekdayOfMonth(date, time.Thursday, 4): // Thanksgiving: 4th Thursday of November
		if month == time.November {
			return true
		}
	}

	return false
}

// IsUKHoliday checks whether the given date is a major UK public (bank) holiday.
// Focused on holidays in England and Wales that can impact the London forex session.
//
// Fixed and variable-date holidays covered:
// - New Year's Day (January 1, or following Monday if weekend)
// - Good Friday (Friday before Easter Sunday)
// - Easter Monday (Monday after Easter Sunday)
// - Early May Bank Holiday (1st Monday of May)
// - Spring Bank Holiday (Last Monday of May)
// - Summer Bank Holiday (Last Monday of August)
// - Christmas Day (December 25, or following Monday if weekend)
// - Boxing Day (December 26, or following weekday if on weekend)
func IsUKHoliday(date time.Time) bool {
	year := date.Year()
	month := date.Month()
	day := date.Day()
	weekday := date.Weekday()

	easterSunday := calculateEasterSunday(year)
	goodFriday := easterSunday.AddDate(0, 0, -2)
	easterMonday := easterSunday.AddDate(0, 0, 1)

	switch {
	// New Year's Day
	case isObservedHoliday(date, time.January, 1):
		return true

	// Good Friday & Easter Monday
	case sameDay(date, goodFriday):
		return true
	case sameDay(date, easterMonday):
		return true

	// Early May Bank Holiday (1st Monday of May)
	case month == time.May && weekday == time.Monday && day <= 7:
		return true

	// Spring Bank Holiday (last Monday of May)
	case month == time.May && weekday == time.Monday && isLastWeekdayOfMonth(date, time.Monday):
		return true

	// Summer Bank Holiday (last Monday of August)
	case month == time.August && weekday == time.Monday && isLastWeekdayOfMonth(date, time.Monday):
		return true

	// Christmas Day & Boxing Day with weekend observance
	case isObservedHoliday(date, time.December, 25):
		return true
	case isObservedHoliday(date, time.December, 26):
		return true
	}

	return false
}

// isNthWeekdayOfMonth returns true if the date is the nth occurrence of the given weekday in its month.
func isNthWeekdayOfMonth(date time.Time, weekday time.Weekday, nth int) bool {
	if date.Weekday() != weekday {
		return false
	}
	// Count how many same weekdays have occurred this month up to this date
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	weekdayCount := 0
	for d := firstOfMonth; d.Month() == date.Month(); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == weekday {
			weekdayCount++
			if d.Day() == date.Day() {
				return weekdayCount == nth
			}
		}
	}
	return false
}

// isLastWeekdayOfMonth returns true if the given date is the last occurrence of the weekday in the month.
func isLastWeekdayOfMonth(date time.Time, weekday time.Weekday) bool {
	if date.Weekday() != weekday {
		return false
	}
	nextWeek := date.AddDate(0, 0, 7)
	return nextWeek.Month() != date.Month()
}

// sameDay checks if two time.Time values are on the same calendar day.
func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

// isObservedHoliday checks if the date matches the actual or observed date of a fixed-date holiday.
// If the holiday falls on a weekend, the following weekday is used.
func isObservedHoliday(date time.Time, month time.Month, day int) bool {
	// Actual holiday
	if date.Month() == month && date.Day() == day {
		return true
	}
	// If holiday falls on Saturday or Sunday, observed on Monday
	holiday := time.Date(date.Year(), month, day, 0, 0, 0, 0, date.Location())
	if holiday.Weekday() == time.Saturday && sameDay(date, holiday.AddDate(0, 0, 2)) {
		return true
	}
	if holiday.Weekday() == time.Sunday && sameDay(date, holiday.AddDate(0, 0, 1)) {
		return true
	}
	return false
}

// calculateEasterSunday returns the date of Easter Sunday for a given year using Anonymous Gregorian algorithm.
func calculateEasterSunday(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Session represents a recurring daily session with opening and closing hours/minutes.
type Session struct {
	name      string
	startHour int
	startMin  int
	endHour   int
	endMin    int
	location  *time.Location
}

// NewSession creates a new trading session with start/end time and timezone.
func NewSession(name string, startHour, startMin, endHour, endMin int, loc *time.Location) *Session {
	return &Session{
		name:      name,
		startHour: startHour,
		startMin:  startMin,
		endHour:   endHour,
		endMin:    endMin,
		location:  loc,
	}
}

// IsOpen returns true if the given time falls within the session (in session's time zone).
func (s *Session) IsOpen(t time.Time) bool {
	localTime := t.In(s.location)

	start := time.Date(localTime.Year(), localTime.Month(), localTime.Day(),
		s.startHour, s.startMin, 0, 0, s.location)

	end := time.Date(localTime.Year(), localTime.Month(), localTime.Day(),
		s.endHour, s.endMin, 0, 0, s.location)

	return !localTime.Before(start) && !localTime.After(end)
}

func (s *Session) String() string {
	return s.name
}

var (
	LondonSession = func() *Session {
		loc, _ := time.LoadLocation("Europe/London")
		return NewSession("London", 8, 0, 17, 0, loc)
	}()

	NYSession = func() *Session {
		loc, _ := time.LoadLocation("America/New_York")
		return NewSession("New York", 9, 0, 17, 0, loc)
	}()
)
