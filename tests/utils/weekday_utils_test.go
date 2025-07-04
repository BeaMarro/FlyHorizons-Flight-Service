package utils_test

import (
	"flyhorizons-flightservice/models/enums"
	"flyhorizons-flightservice/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type WeekDayUtilsTest struct {
}

// Setup
func setupWeekdayUtils() utils.WeekdayUtils {
	return utils.WeekdayUtils{}
}

func TestConvertDateToWeekdayReturnsDay(t *testing.T) {
	// Arrange
	date := time.Date(2025, time.March, 14, 0, 0, 0, 0, time.Local)
	weekdayUtils := setupWeekdayUtils()
	// Act
	day := weekdayUtils.ConvertToWeekDay(date)
	// Assert
	assert.Equal(t, day, enums.Friday)
}

func TestListContainingDayContainsReturnsTrue(t *testing.T) {
	// Arrange
	days := []enums.Day{enums.Monday, enums.Friday}
	targetDay := enums.Friday
	weekdayUtils := setupWeekdayUtils()
	// Act
	containsDay := weekdayUtils.ContainsDay(days, targetDay)
	// Assert
	assert.True(t, containsDay)
}

func TestListNotContainingDayContainsReturnsFalse(t *testing.T) {
	// Arrange
	days := []enums.Day{enums.Monday, enums.Friday}
	targetDay := enums.Tuesday
	weekdayUtils := setupWeekdayUtils()
	// Act
	containsDay := weekdayUtils.ContainsDay(days, targetDay)
	// Assert
	assert.False(t, containsDay)
}

func TestConvertJSONToDaysReturnsListDay(t *testing.T) {
	// Arrange
	daysJSON := "[1, 3, 5]"
	weekdayUtils := setupWeekdayUtils()
	expectedDays := []enums.Day{enums.Monday, enums.Wednesday, enums.Friday}
	// Act
	days, err := weekdayUtils.ConvertJSONToDays(daysJSON)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDays, days)
}

func TestConvertDaysToJSONReturnsJSON(t *testing.T) {
	// Arrange
	days := []enums.Day{enums.Monday, enums.Wednesday, enums.Friday}
	weekdayUtils := setupWeekdayUtils()
	expectedJSON := "[1,3,5]"
	// Act
	JSON, err := weekdayUtils.ConvertDaysToJSON(days)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, JSON)
}
