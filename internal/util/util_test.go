package util

import (
	"testing"
	"time"
)

//revive:disable:cyclomatic
func TestDaysElapsedAbsolute(t *testing.T) {
	sameTimeSameDayA, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}
	sameTimeSameDayB, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	differentTimeSameDayA, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.48Z")
	if err != nil {
		t.Fatal(err)
	}
	differentTimeSameDayB, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	sameTimeDifferentDayA, err := time.Parse(time.RFC3339Nano, "2018-10-29T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}
	sameTimeDifferentDayB, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	sameTimeDifferentDayASeveralDays, err := time.Parse(time.RFC3339Nano, "2018-10-31T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}
	sameTimeDifferentDayBSeveralDays, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	differentTimeDifferentDayA, err := time.Parse(time.RFC3339Nano, "2018-10-29T04:44:08.48Z")
	if err != nil {
		t.Fatal(err)
	}
	differentTimeDifferentDayB, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	differentTimeDifferentDayAWithinOneNano, err := time.Parse(time.RFC3339Nano, "2018-10-29T04:44:08.46Z")
	if err != nil {
		t.Fatal(err)
	}
	differentTimeDifferentDayBWithinOneNano, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}

	differentTimeDifferentDayAWithinOneNanoOppositeOrder, err := time.Parse(time.RFC3339Nano, "2018-10-28T04:44:08.47Z")
	if err != nil {
		t.Fatal(err)
	}
	differentTimeDifferentDayBWithinOneNanoOppositeOrder, err := time.Parse(time.RFC3339Nano, "2018-10-29T04:44:08.46Z")
	if err != nil {
		t.Fatal(err)
	}

	type input struct {
		dateA time.Time
		dateB time.Time
	}

	var data = []struct {
		desc     string
		input    input
		expected int
	}{
		{
			desc: "same time same day",
			input: input{
				dateA: sameTimeSameDayA,
				dateB: sameTimeSameDayB,
			},
			expected: 1,
		},
		{
			desc: "different time same day",
			input: input{
				dateA: differentTimeSameDayA,
				dateB: differentTimeSameDayB,
			},
			expected: 1,
		},
		{
			desc: "same time different day",
			input: input{
				dateA: sameTimeDifferentDayA,
				dateB: sameTimeDifferentDayB,
			},
			expected: 2,
		},
		{
			desc: "same time different day, several days",
			input: input{
				dateA: sameTimeDifferentDayASeveralDays,
				dateB: sameTimeDifferentDayBSeveralDays,
			},
			expected: 4,
		},
		{
			desc: "different time different day",
			input: input{
				dateA: differentTimeDifferentDayA,
				dateB: differentTimeDifferentDayB,
			},
			expected: 2,
		},
		{
			desc: "different time different day within on nano of ticking over",
			input: input{
				dateA: differentTimeDifferentDayAWithinOneNano,
				dateB: differentTimeDifferentDayBWithinOneNano,
			},
			expected: 1,
		},
		{
			desc: "different time different day within on nano of ticking over, opposite order",
			input: input{
				dateA: differentTimeDifferentDayAWithinOneNanoOppositeOrder,
				dateB: differentTimeDifferentDayBWithinOneNanoOppositeOrder,
			},
			expected: 1,
		},
	}

	for i, d := range data {
		result := DaysElapsedAbsolute(d.input.dateA, d.input.dateB)

		if result != d.expected {
			t.Errorf("Description=%s, At=%d. Unexpected result content. Result: %v || Expected: %v", d.desc, i, result, Print(d.expected))
		}
	}
}

//revive:enable:cyclomatic
