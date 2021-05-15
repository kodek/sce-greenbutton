package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsHoliday_ChristmasIsTrue(t *testing.T) {
	date := time.Date(2021, 12, 25, 12, 23, 35, 0, time.UTC)

	got := IsHoliday(date)

	assert.True(t, got)
}
func TestIsHoliday_WeekBeforeChristmasIsFalse(t *testing.T) {
	date := time.Date(2021, 12, 25, 12, 23, 35, 0, time.UTC).Add(-7 * 24 * time.Hour)

	got := IsHoliday(date)

	assert.False(t, got)
}
