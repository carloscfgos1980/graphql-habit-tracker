package utils

import (
	"time"

	"github.com/carloscfgos1980/graphql-habit-tracker/internal/models"
)

// 456fe4f65wef4we65f4s65f4v65fs4v = Morning Walk
// [J25, J24, J23, J21, J20, J19, J18, J17, J15, J14, J13, J12, J11, J10, J09]
func CalculateCurrentStreak(logs []*models.HabitLog) int {
	if len(logs) == 0 {
		return 0
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	firstLogDate := logs[0].CompletedDate

	if !firstLogDate.Equal(today) && !firstLogDate.Equal(yesterday) {
		return 0
	}

	streak := 1

	prevDate := firstLogDate

	for i := 1; i < len(logs); i++ {
		logDate := logs[i].CompletedDate
		expectedDate := prevDate.AddDate(0, 0, -1)

		if logDate.Equal(expectedDate) {
			streak++

			prevDate = logDate
		} else {
			break
		}
	}

	return streak
}

// 456fe4f65wef4we65f4s65f4v65fs4v = Morning Walk
// Original Slice - DESC - Newest to Oldest
// [J25, J24, J23, J21, J20, J19, J18, J17, J15, J14, J13, J12, J11, J10, J09]

// i=0, j14 ->>>>> 0<14
// [J09, J24, J23, J21, J20, J19, J18, J17, J15, J14, J13, J12, J11, J10, J25]
// i=1, j13 ->>>>> 1<13
// [J09, J10, J23, J21, J20, J19, J18, J17, J15, J14, J13, J12, J11, J24, J25]
// i=2, j12 ->>>>> 2<12
// [J09, J10, J11, J21, J20, J19, J18, J17, J15, J14, J13, J12, J23, J24, J25]
// i=3, j11 ->>>>> 3<11
// [J09, J10, J11, J12, J20, J19, J18, J17, J15, J14, J13, J21, J23, J24, J25]
// i=4, j10 ->>>>> 4<10
// [J09, J10, J11, J12, J13, J19, J18, J17, J15, J14, J20, J21, J23, J24, J25]
// i=5, j9 ->>>>> 5<9
// [J09, J10, J11, J12, J13, J14, J18, J17, J15, J19, J20, J21, J23, J24, J25]
// i=6, j8 ->>>>> 6<8
// [J09, J10, J11, J12, J13, J14, J15, J17, J18, J19, J20, J21, J23, J24, J25]

// i=7, j7 ->>>>> 7<7 STOP

// Reversed Slice - ASC - Oldest to Newest
// [J09, J10, J11, J12, J13, J14, J15, J17, J18, J19, J20, J21, J23, J24, J25]
func CalculateLongestStreak(logs []*models.HabitLog) int {
	if len(logs) == 0 {
		return 0
	}

	// i, j
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	currentRun := 1
	longestRun := 1

	prevDate := logs[0].CompletedDate.UTC().Truncate(24 * time.Hour)

	// [J09, J10, J11, J12, J13, J14, J15, J17, J18, J19, J20, J21, J23, J24, J25]
	for i := 1; i < len(logs); i++ {
		// J17
		logDate := logs[i].CompletedDate.UTC().Truncate(24 * time.Hour)
		// J15 + 1 = J16
		expectedDate := prevDate.AddDate(0, 0, 1)

		if logDate.Equal(expectedDate) {
			currentRun++

			if currentRun > longestRun {
				longestRun = currentRun
			}
		} else {
			currentRun = 1
		}

		prevDate = logDate
	}

	return longestRun
}
