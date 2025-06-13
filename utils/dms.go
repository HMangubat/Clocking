package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ParseDMS converts DMS strings to decimal degrees.
// Supports formats like:
// - "12°36′15.47″ N"
// - "14:09:12.42 N"
// - "121:15:58.30 E"
func ParseDMS(dms string) (float64, error) {
	// Trim and normalize whitespace
	dms = strings.TrimSpace(dms)

	// Try colon-separated format first: "14:09:12.42 N"
	reColon := regexp.MustCompile(`^(\d+):(\d+):([\d.]+)\s*([NSEW])$`)
	if matches := reColon.FindStringSubmatch(dms); len(matches) == 5 {
		return parseDMSValues(matches[1], matches[2], matches[3], matches[4])
	}

	// Fallback to unicode DMS format: "12°36′15.47″ N"
	reUnicode := regexp.MustCompile(`^(\d+)°(\d+)′([\d.]+)″\s*([NSEW])$`)
	if matches := reUnicode.FindStringSubmatch(dms); len(matches) == 5 {
		return parseDMSValues(matches[1], matches[2], matches[3], matches[4])
	}

	return 0, errors.New("invalid DMS format")
}

// parseDMSValues performs conversion to decimal degrees.
func parseDMSValues(degStr, minStr, secStr, dir string) (float64, error) {
	deg, err1 := strconv.ParseFloat(degStr, 64)
	min, err2 := strconv.ParseFloat(minStr, 64)
	sec, err3 := strconv.ParseFloat(secStr, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, errors.New("invalid number in DMS")
	}

	decimal := deg + min/60 + sec/3600
	if dir == "S" || dir == "W" {
		decimal = -decimal
	}
	return decimal, nil
}
