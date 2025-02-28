package iprange_test

import (
	"gitlab.com/slon/shad-go/iprange"
	"testing"
)

func initIPPatterns() []string {
	return []string{
		"10.0.0.5-10",
		"10.0.0.1",
		"192.168.10.0/24",
		"192.168.1.*",
	}
}

func FuzzParseList(f *testing.F) {
	presetPatterns := initIPPatterns()
	feedPatternsToFuzzer(f, presetPatterns)

	f.Fuzz(func(t *testing.T, ipRangeString string) {
		executeParseWithPanicCheck(t, ipRangeString)
	})
}

func feedPatternsToFuzzer(f *testing.F, patterns []string) {
	for _, pattern := range patterns {
		addPatternToFuzzer(f, pattern)
	}
}

func addPatternToFuzzer(f *testing.F, pattern string) {
	f.Add(pattern)
}

func executeParseWithPanicCheck(t *testing.T, ipRangeString string) {
	defer recoverAndLogIfPanic(t, ipRangeString)
	tryParseAndLogResult(t, ipRangeString)
}

func recoverAndLogIfPanic(t *testing.T, ipRangeString string) {
	if r := recover(); r != nil {
		logPanicOccurred(t, ipRangeString, r)
	}
}

func tryParseAndLogResult(t *testing.T, ipRangeString string) {
	if result, err := iprange.Parse(ipRangeString); err != nil {
		reportParseError(t, ipRangeString, err)
	} else {
		reportParseSuccess(t, ipRangeString, result)
	}
}

func logPanicOccurred(t *testing.T, input string, panicDetail interface{}) {
	t.Errorf("Обнаружена паника при парсинге %q: %v", input, panicDetail)
}

func reportParseError(t *testing.T, input string, err error) {
	t.Logf("Ошибка при обработке %q: %v", input, err)
}

func reportParseSuccess(t *testing.T, input string, result interface{}) {
	t.Logf("Успешный результат парсинга %q: %v", input, result)
}
