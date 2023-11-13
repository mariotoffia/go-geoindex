package geoindex_test

import (
	"math"
	"testing"

	geoindex "github.com/mariotoffia/go-geoindex"
	"github.com/stretchr/testify/assert"
)

func TestMovePointShortDistancesSmallDistance(t *testing.T) {
	from := geoindex.GeoPoint{Plat: 51.0, Plon: 0.0} // Some reference point
	direction := geoindex.North
	distanceMeters := 10000.0 // 10 km
	expectedLat := 51.09      // Approximate expected latitude
	name := "TestPoint"

	result := geoindex.MovePointShortDistances(&from, distanceMeters, direction, name)
	diff := math.Abs(result.Plat - expectedLat)

	assert.True(t, diff < 0.01, "MovePointShortDistances failed for small distance: expected latitude around %v, got %v", expectedLat, result.Plat)
}

func TestMovePointShortDistancesBoundaryDistance(t *testing.T) {
	from := geoindex.GeoPoint{Plat: 51.0, Plon: 0.0} // Reference point
	direction := geoindex.North
	distanceMeters := 100000.0 // 100 km, near the upper limit of the method's accuracy
	expectedLat := 51.9        // Approximate expected latitude
	name := "TestPointBoundary"

	result := geoindex.MovePointShortDistances(&from, distanceMeters, direction, name)
	diff := math.Abs(result.Plat - expectedLat)

	assert.True(t, diff < 0.01, "MovePointShortDistances failed for boundary distance: expected latitude around %v, got %v", expectedLat, result.Plat)
}

func TestMovePointShortDistancesBoundaryDistance2(t *testing.T) {
	from := geoindex.GeoPoint{Plat: 51.0, Plon: 0.0} // Reference point
	direction := geoindex.North
	distanceMeters := 500000.0 // 500 km, near the upper limit of the method's accuracy
	expectedLat := 55.5        // Approximate expected latitude
	name := "TestPointBoundary"

	result := geoindex.MovePointShortDistances(&from, distanceMeters, direction, name)
	diff := math.Abs(result.Plat - expectedLat)

	assert.True(t, diff < 0.01, "MovePointShortDistances failed for boundary distance: expected latitude around %v, got %v", expectedLat, result.Plat)
}

func TestMovePointLongDistancesLargeDistance(t *testing.T) {
	//t.Skip("not implemented (need to verify that it works)")

	from := geoindex.GeoPoint{Plat: 59.32813965282047, Plon: 18.065817040989696} // Reference point
	direction := geoindex.West
	distanceMeters := float64(6319 * 1000) // 6319 km, need to switch to more accurate calculation
	expectedLat := 59.993                  // Approximate expected latitude
	name := "TestPointBoundary"

	// 40.725850378291774, -73.93621697956644
	result := geoindex.MovePointShortDistances(&from, distanceMeters, direction, name)
	//result, err := geoindex.MovePointLongDistances(&from, distanceMeters, direction, name)
	//require.NoError(t, err)

	diff := math.Abs(result.Plat - expectedLat)

	assert.True(t, diff < 0.01, "MovePointShortDistances failed for boundary distance: expected latitude around %v, got %v", expectedLat, result.Plat)
}
