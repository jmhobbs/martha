package synthetic

import (
	"testing"
	"time"

	"github.com/jmhobbs/martha/internal/configuration"
)

func Test_New(t *testing.T) {
	_, err := New(configuration.Synthetic{
		Interval: "invalid",
	})
	if err == nil {
		t.Errorf("Expected error for invalid Interval, got nil")
	}

	_, err = New(configuration.Synthetic{
		Interval: "15s",
		Duration: "invalid",
	})
	if err == nil {
		t.Errorf("Expected error for invalid Interval, got nil")
	}
}

func Test_SyntheticReadForTime(t *testing.T) {
	now := time.Now()

	// Active for 30 seconds, every 15 second.
	s, err := New(configuration.Synthetic{
		Interval: "15s",
		Duration: "30s",
		Effect: configuration.CONTROL_EFFECT_TYPE_HUMIDITY,
		Values: configuration.SyntheticValues{
			Idle: 95.0,
			Active: 90.0,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Last Active state started NOW
	s.lastActive = now

	/*
	Activates every 45s (interval + duration)

	TS	Last Active		Active?
	 0	0s ago				Yes
	15	15s ago				Yes
	30	30s ago				No
	31	31s ago				No
	45	45s ago				Yes
	55  10s ago				Yes
	75	30s ago				No
	76	31s ago				No
	*/
	expectations := []struct{
		Offset int
		ExpectActive bool
	} {
		{ 0, true },
		{ 15, true },
		{ 30, false },
		{ 31, false },
		{ 45, true },
		{ 55, true },
		{ 75, false },
		{ 76, false },
	}

	for _, expectation := range expectations {
		thisNow := now.Add(time.Second * time.Duration(expectation.Offset))
		reading := s.readForTime(thisNow)
		expected := s.idleValue
		if expectation.ExpectActive {
			expected = s.activeValue
		}
		if *reading.Humidity != expected {
			t.Errorf(
				"Incorrect reading for t=%d\nexpected: %v\n  actual: %v",
				expectation.Offset,
				expected,
				*reading.Humidity,
			)
		}		
	}
}