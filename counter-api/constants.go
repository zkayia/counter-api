package main

import (
	"errors"
)

type CounterUpdate struct {
	time      int
	counter   string
	value     int
	operation Operation
	amount    int
}

type Duration int

const (
	Infinity Duration = iota
	Second
	Minute
	Hour
	Day
	Month
	Year
	Decade
	Century
)

func (duration Duration) toString() string {
	return [...]string{
		"infinity",
		"second",
		"minute",
		"hour",
		"day",
		"month",
		"year",
		"decade",
		"century",
	}[duration]
}

func durationFromString(value string) (Duration, error) {
	switch value {
	case "second":
		return Second, nil
	case "minute":
		return Minute, nil
	case "hour":
		return Hour, nil
	case "day":
		return Day, nil
	case "month":
		return Month, nil
	case "year":
		return Year, nil
	case "decade":
		return Decade, nil
	case "century":
		return Century, nil
	}
	return 0, errors.New("Invalid duration string: `" + value + "`.")
}

type Operation int

const (
	Add Operation = iota
	Sub
	Set
)

var Operations = [...]Operation{
	Add,
	Sub,
	Set,
}

func (operation Operation) toString() string {
	return [...]string{
		"add",
		"sub",
		"set",
	}[operation]
}
