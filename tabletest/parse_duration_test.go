package tabletest

import (
	"math/rand"
	"testing"
	"time"
)

var parseDurationTests = []struct {
	input    string
	duration time.Duration
	ok       bool
}{
	{"0", 0, true},
	{"2s", 2 * time.Second, true},
	{"15m", 15 * time.Minute, true},
	{"8h", 8 * time.Hour, true},
	{"3600ms", 3600 * time.Millisecond, true},
	{"1200ns", 1200 * time.Nanosecond, true},
	{"1000us", 1000 * time.Microsecond, true},
	{"750µs", 750 * time.Microsecond, true}, // Microsecond symbol

	{"2h30m", 2*time.Hour + 30*time.Minute, true},
	{"45m15s", 45*time.Minute + 15*time.Second, true},
	{"1h20m5s200ms", 1*time.Hour + 20*time.Minute + 5*time.Second + 200*time.Millisecond, true},
	{"2h15.75m", 2*time.Hour + 15*time.Minute + 45*time.Second, true},
	{"5m0.75s", 5*time.Minute + 750*time.Millisecond, true},

	{"-10s", -10 * time.Second, true},
	{"+5m", 5 * time.Minute, true},
	{"-0", 0, true},
	{"+0", 0, true},

	{"3.5h", 3*time.Hour + 30*time.Minute, true},
	{"2.25m", 2*time.Minute + 15*time.Second, true},
	{"1.5s", 1*time.Second + 500*time.Millisecond, true},
	{"0.25s", 250 * time.Millisecond, true},
	{"7.75m", 7*time.Minute + 45*time.Second, true},

	{"1h2m3s", 1*time.Hour + 2*time.Minute + 3*time.Second, true},
	{"4h30m15.5s", 4*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond, true},
	{"-3m10s", -(3*time.Minute + 10*time.Second), true},
	{"10s20ms5µs", 10*time.Second + 20*time.Millisecond + 5*time.Microsecond, true},
	{"1h2m3s4ms5ns6µs", 1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond + 5*time.Nanosecond + 6*time.Microsecond, true},

	{"92233720368ns", 92233720368 * time.Nanosecond, true},
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond, true},
	{"0.000000001s", 1 * time.Nanosecond, true},
	{"72000m", 72000 * time.Minute, true},

	{"0.3333333333333333333h", 20 * time.Minute, true},  // https://golang.org/issue/6617
	{"0.100000000000000000000h", 6 * time.Minute, true}, // https://golang.org/issue/15011

	{"0.9223372036854775807h", 55*time.Minute + 20*time.Second + 413933267*time.Nanosecond, true},

	{"", 0, false},
	{"5", 0, false},         // No unit
	{"2h-5m", 0, false},     // Incorrect negative
	{"2..5s", 0, false},     // Double dot
	{"1s1", 0, false},       // No unit for second number
	{"m10", 0, false},       // Invalid ordering
	{"s", 0, false},         // Unit only
	{"10minutes", 0, false}, // Incorrect unit name
	{".", 0, false},         // Single dot
	{"s.", 0, false},        // Unit then dot
	{"-.5s", -500 * time.Millisecond, true},
	{"+.5s", 500 * time.Millisecond, true},

	{"3000000h", 0, false},                  // Overflow
	{"9223372036854775808ns", 0, false},     // Overflow by one nanosecond
	{"9223372036854ms775μs808ns", 0, false}, // Large mixed overflow
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		res, err := ParseDuration(tc.input)
		if tc.ok && (err != nil || res != tc.duration) {
			t.Errorf("ParseDuration(%q) = (%v, %v), want (%v, nil)", tc.input, res, err, tc.duration)
		} else if !tc.ok && err == nil {
			t.Errorf("ParseDuration(%q) = (_, nil), want (_, non-nil)", tc.input)
		}
	}
}

func TestParseDurationRoundTrip(t *testing.T) {
	for i := 0; i < 100; i++ {
		res0 := time.Duration(rand.Int31()) * time.Millisecond
		s := res0.String()
		res1, err := ParseDuration(s)
		if err != nil || res0 != res1 {
			t.Errorf("round-trip failed: %d => %q => %d, %v", res0, s, res1, err)
		}
	}
}
