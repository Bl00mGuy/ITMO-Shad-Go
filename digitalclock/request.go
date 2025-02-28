package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type ClockRequest struct {
	Scale int
	Time  string
}

func ParseClockRequest(request *http.Request) (*ClockRequest, error) {
	scale, err := parseScale(request)
	if err != nil {
		return nil, err
	}

	timeParameter, err := parseTime(request)
	if err != nil {
		return nil, err
	}

	return &ClockRequest{
		Scale: scale,
		Time:  timeParameter,
	}, nil
}

func parseScale(request *http.Request) (int, error) {
	query := request.URL.Query()
	if k, exists := query["k"]; exists && len(k) > 0 {
		scale, err := strconv.Atoi(k[0])
		if err != nil || scale < 1 || scale > 30 {
			return 0, errors.New("invalid k")
		}
		return scale, nil
	}
	return 1, nil
}

func parseTime(request *http.Request) (string, error) {
	query := request.URL.Query()
	timeParameter := query.Get("time")

	if timeParameter == "" {
		return getCurrentTime(), nil
	}

	if !isValidTimeFormat(timeParameter) {
		return "", errors.New("invalid time")
	}

	return timeParameter, nil
}

func isValidTimeFormat(timeString string) bool {
	timePattern := `^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`
	match, _ := regexp.MatchString(timePattern, timeString)
	return match
}

func getCurrentTime() string {
	now := time.Now()
	return fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
}
