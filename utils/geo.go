package utils

import "math"

// RoundTo6Decimals rounds a float to 6 decimal places.
func RoundTo6Decimals(f float64) float64 {
	return math.Round(f*1_000_000) / 1_000_000
}

// RoundTo3Decimals rounds a float to 3 decimal places.
func RoundTo3Decimals(f float64) float64 {
	return math.Round(f*1000) / 1000
}

// HaversineDistance returns the distance between two lat/long points in kilometers.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }

	dlat := toRad(lat2 - lat1)
	dlon := toRad(lon2 - lon1)
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
