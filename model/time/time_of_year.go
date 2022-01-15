package time

type TimeOfYear struct {
	Month uint8
	Day   uint8
	Hour  uint8
}

func (t *TimeOfYear) Matches(c *CalendarType) bool {
	return t.Month == c.Month && t.Day == c.Day && t.Hour == c.Hour
}
