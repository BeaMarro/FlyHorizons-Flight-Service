package utils

import (
	"encoding/json"
	"flyhorizons-flightservice/models/enums"
	"time"
)

type WeekdayUtils struct{}

func (utils WeekdayUtils) getDayOfWeek(date time.Time) string {
	return date.Weekday().String()
}

func (utils WeekdayUtils) convertToDay(day string) enums.Day {
	switch day {
	case "Monday":
		return enums.Monday
	case "Tuesday":
		return enums.Tuesday
	case "Wednesday":
		return enums.Wednesday
	case "Thursday":
		return enums.Thursday
	case "Friday":
		return enums.Friday
	case "Saturday":
		return enums.Saturday
	case "Sunday":
		return enums.Sunday
	default:
		return 0 // Invalid day
	}
}

func (utils WeekdayUtils) ConvertToWeekDay(date time.Time) enums.Day {
	dayString := utils.getDayOfWeek(date)
	return utils.convertToDay(dayString)
}

func (utils WeekdayUtils) ContainsDay(days []enums.Day, targetDay enums.Day) bool {
	for _, day := range days {
		if enums.Day(day) == targetDay {
			return true
		}
	}
	return false
}

func (utils WeekdayUtils) ConvertJSONToDays(jsonString string) ([]enums.Day, error) {
	var dayInts []int
	err := json.Unmarshal([]byte(jsonString), &dayInts)
	if err != nil {
		return nil, err
	}

	days := make([]enums.Day, len(dayInts))
	for i, dayInt := range dayInts {
		days[i] = enums.Day(dayInt)
	}

	return days, nil
}

func (utils WeekdayUtils) ConvertDaysToJSON(days []enums.Day) (string, error) {
	dayInts := make([]int, len(days))
	for i, day := range days {
		dayInts[i] = int(day)
	}

	jsonBytes, err := json.Marshal(dayInts)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
