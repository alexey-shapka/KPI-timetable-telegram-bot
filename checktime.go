package main

import (
	"time"
)

func checkWeek() int{
	_, week := time.Now().ISOWeek()
	switch week%2{
	case 0:
		week = 2
	case 1:
		week = 1
	}

	return week
}

func checkDay() int{
	return int(time.Now().Weekday())
}