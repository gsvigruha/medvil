package controller

import (
//"fmt"
)

type CalendarType struct {
	Year  uint16
	Month uint8
	Day   uint8
	Hour  uint8
}

const (
	Spring uint8 = 0
	Summer uint8 = 1
	Autumn uint8 = 2
	Winter uint8 = 3
)

var Calendar = &CalendarType{
	Year:  1000,
	Month: 3,
	Day:   1,
	Hour:  0,
}

func (c *CalendarType) Tick() {
	c.Hour++
	if c.Hour == 24 {
		c.Hour = 0
		c.Day++
	}
	if c.Day == 31 {
		c.Day = 1
		c.Month++
	}
	if c.Month == 13 {
		c.Month = 1
		c.Year++
	}
}

func (c *CalendarType) Season() uint8 {
	if c.Month >= 3 && c.Month <= 5 {
		return Spring
	}
	if c.Month >= 6 && c.Month <= 8 {
		return Summer
	}
	if c.Month >= 9 && c.Month <= 11 {
		return Autumn
	}
	return Winter
}
