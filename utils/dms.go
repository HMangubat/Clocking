package utils

import (
	"errors"
	"regexp"
	"strconv"
)

// ParseDMS converts a DMS string like "12°36′15.47″ N" to decimal degrees.
func ParseDMS(dms string) (float64, error) {
	re := regexp.MustCompile(`(\d+)°(\d+)′([\d.]+)″\s*([NSEW])`)
	matches := re.FindStringSubmatch(dms)
	if len(matches) != 5 {
		return 0, errors.New("invalid DMS format")
	}
	deg, _ := strconv.ParseFloat(matches[1], 64)
	min, _ := strconv.ParseFloat(matches[2], 64)
	sec, _ := strconv.ParseFloat(matches[3], 64)
	dir := matches[4]

	decimal := deg + min/60 + sec/3600
	if dir == "S" || dir == "W" {
		decimal = -decimal
	}
	return decimal, nil
}
