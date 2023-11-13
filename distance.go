package geoindex

import (
	"errors"
	"math"
)

// TODO: Need to investigate https://www.movable-type.co.uk/scripts/latlong.html for calculations!

const (
	er                = 6378137           // Semi-major axis of the Earth (meters) - a more precise measurement than EarthRadius
	flattening        = 1 / 298.257223563 // Flattening of the Earth
	bearingAdjustment = 0.0001            // Adjustment to bearing in radians for each iteration
	maxIterations     = 200               // Maximum number of iterations for convergence
)

// MovePointShortDistances will move point _from_ to a new point _distanceMeters_ away in the given _direction_.
//
// The new point is returned and the _from_ is untouched.
//
// Depending on the distance, this function will use either MovePointShortDistances or MovePointLongDistances.
// Where the latter is computational intensive and may even fail.
func MovePoint(from Point, distanceMeters float64, direction Direction, name string) (*GeoPoint, error) {
	if distanceMeters < 500000 /* 500km */ {
		return MovePointShortDistances(from, distanceMeters, direction, name), nil
	}

	return nil, errors.New("not implemented (need to verify that it works)")
	//return MovePointLongDistances(from, distanceMeters, direction, name)
}

// MovePointShortDistances will move point _from_ to a new point _distanceMeters_ away in the given _direction_.
//
// The new point is returned and the _from_ is untouched.
//
// CAUTION: This function assumes that the earth is a perfect sphere and hence is not accurate for large distances.
// Use it for distances up to 500km.
//
// See http://en.wikipedia.org/wiki/Geographical_distance#Spherical_Earth_projected_to_a_plane
func MovePointShortDistances(from Point, distanceMeters float64, direction Direction, name string) *GeoPoint {
	// Convert distance to radians
	deltaDistance := distanceMeters / float64(earthRadius)
	bearing := direction.ToRadians()

	lat1 := toRadians(from.Lat())
	lon1 := toRadians(from.Lon())

	// Calculate the new lat and long
	newLat := math.Asin(math.Sin(lat1)*math.Cos(deltaDistance) +
		math.Cos(lat1)*math.Sin(deltaDistance)*math.Cos(bearing))

	newLon := lon1 + math.Atan2(math.Sin(bearing)*math.Sin(deltaDistance)*math.Cos(lat1),
		math.Cos(deltaDistance)-math.Sin(lat1)*math.Sin(newLat))

	return &GeoPoint{
		Pid:  name,
		Plat: toDegrees(newLat),
		Plon: toDegrees(newLon),
	}
}

// MovePointLongDistances calculates new latitude and longitude using the Vincenty formula with fallback for non-convergence.
//
// The new point is returned and the _from_ is untouched.
//
// CAUTION: This function assumes that the earth is an ellipsoid and hence is more accurate for large distances but is much
// more expensive than MovePointShortDistances in terms of computation power. Sometimes the algorithm does not converge
// and in that case it will return an error.
//
// See http://en.wikipedia.org/wiki/Vincenty%27s_formulae
func MovePointLongDistances(pt Point, distanceMeters float64, direction Direction, name string) (*GeoPoint, error) {
	// Initial bearing in radians
	initialBearing := direction.ToRadians()
	bearing := initialBearing

	for i := 0; i < maxIterations; i++ {
		finalLat, finalLon, err := calculateVincenty(toRadians(pt.Lat()), toRadians(pt.Lon()), distanceMeters, bearing)
		if err == nil {
			return &GeoPoint{
				Pid:  name,
				Plat: finalLat,
				Plon: finalLon,
			}, nil
		}

		// Adjust bearing slightly
		bearing = initialBearing + float64(i+1)*bearingAdjustment
	}

	return nil, errors.New("failed to converge")
}

func calculateVincenty(lat, lon, distanceMeters, bearing float64) (float64, float64, error) {
	b := er * (1 - flattening) // Semi-minor axis
	tanU1 := (1 - flattening) * math.Tan(lat)
	cosU1 := 1 / math.Sqrt((1 + tanU1*tanU1))
	sinU1 := tanU1 * cosU1
	cosBearing := math.Cos(bearing)
	sinBearing := math.Sin(bearing)

	// Calculations
	u1 := math.Atan((1 - flattening) * math.Tan(lat))
	sigma1 := math.Atan2(math.Tan(u1), math.Cos(bearing))
	sinAlpha := math.Cos(u1) * math.Sin(bearing)
	cosSqAlpha := 1 - sinAlpha*sinAlpha
	uSq := cosSqAlpha * (er*er - b*b) / (b * b)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))

	sigma := distanceMeters / (b * A)
	var sigmaP float64 = 2 * math.Pi
	var cos2SigmaM, sinSigma, cosSigma float64

	for math.Abs(sigma-sigmaP) > 1e-12 {
		cos2SigmaM = math.Cos(2*sigma1 + sigma)
		sinSigma = math.Sin(sigma)
		cosSigma = math.Cos(sigma)
		deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-
			B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))
		sigmaP = sigma
		sigma = distanceMeters/(b*A) + deltaSigma
	}

	tmp := sinU1*sinSigma - cosU1*cosSigma*cosBearing
	lat2 := math.Atan2(sinU1*cosSigma+cosU1*sinSigma*cosBearing, (1-flattening)*math.Sqrt(sinAlpha*sinAlpha+tmp*tmp))
	lambda := math.Atan2(sinSigma*sinBearing, cosU1*cosSigma-sinU1*sinSigma*cosBearing)
	C := flattening / 16 * cosSqAlpha * (4 + flattening*(4-3*cosSqAlpha))
	L := lambda - (1-C)*flattening*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))

	// Final calculations for latitude and longitude
	finalLat := toDegrees(lat2)
	finalLon := toDegrees(lon + L)

	// Check for convergence
	if math.Abs(sigma-sigmaP) <= 1e-12 {
		return finalLat, finalLon, nil
	}

	return 0, 0, errors.New("iteration did not converge")
}
