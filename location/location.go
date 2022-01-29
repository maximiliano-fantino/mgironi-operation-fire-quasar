package location

import (
	"fmt"

	"log"

	"math"

	"github.com/mgironi/operation-fire-quasar/model"

	"github.com/mgironi/operation-fire-quasar/store"
	"github.com/montanaflynn/stats"
)

// input: distance to the transmitter recieved on each satlelite
// output: the coordinates 'x' and 'y' of the message emiter
func GetLocation(distances ...float32) (x, y float32) {
	if len(distances) < 3 {
		log.Printf("Is not possible to compelete calculations. There is not enough distances values as parameter to determine location. At least 3 distances are required ")
	}
	var err error
	x, y, err = CalculateLocation(distances)
	if err != nil {
		log.Printf("Is no possible to compelete calculations. %s", err.Error())
	}
	return x, y
}

// Calculates coordinates location.
// The array should have 3 ordered distances to the known coordinates.
// input: Recieves distances array to a known coordinates.
// output: Returns X and Y coordinates of the calculated location and an error in case calculation couldn't be done.
func CalculateLocation(distances []float32) (x, y float32, err error) {

	// gets reference points coordinates
	pointsCoordinates := GetKnownReferenceCoordinates()

	// checks if distances has same amount of elements that the refences points coordiantes.
	if len(distances) != len(pointsCoordinates) {
		return 0, 0, fmt.Errorf("distances and Points coordinates has diferent sizes. Distances: %d, PointsCoord: %d", len(distances), len(pointsCoordinates))
	}

	// calculates location with the distances to the points coordinates using trilateration math method
	x, y = CalculateLocationByTrilateration(distances, pointsCoordinates)

	// checks if calculated coorinates match with given distances, getting ratioError
	ratioErr, err := ChecksDistancesToCoordinate(distances, pointsCoordinates, x, y)
	if err != nil {
		log.Print(err)
		return 0, 0, err
	}

	// checks if ratioError is acceptable with float comparission tolerance
	acceptedRatio := model.FLOAT_COMPARISION_TOLERANCE
	if ratioErr > acceptedRatio {
		log.Printf("WARN ratio error exceeds aceptable level of %.4f. Ratio ~ %.4f", acceptedRatio, ratioErr)
	}
	return x, y, nil
}

// Routput: the kwnown reference coordinates.
func GetKnownReferenceCoordinates() (points []model.Point) {
	if len(store.Satelites) == 0 {
		store.InitializeSatelitesInfo()
	}
	points = make([]model.Point, len(store.Satelites))
	for i, satelite := range store.Satelites {
		points[i] = satelite.Location
	}
	return points
}

// Checks if the X, Y coordinates distance to each pointsCoordinates matchs with the given distances.
// input: distances, points coordinates and 'x','y' calculated coordinates to check.
// output: the median errorRatio calculated (0: no error, interval [0,1]: percent error)
// error1: if detects arrays length diferences (between distances and points coordinates)
// error2: an internal calculation error.
func ChecksDistancesToCoordinate(distances []float32, pointsCoordinates []model.Point, x, y float32) (errorRatio float64, err error) {
	// checks arrays length, they must be equals
	if len(distances) != len(pointsCoordinates) {
		return 0, fmt.Errorf("can't check distances with coordinate. Distances and Points coordinates has diferent sizes. Distances: %d, PointsCoord: %d", len(distances), len(pointsCoordinates))
	}

	// calculate differences ratios form given distances to calculated distances
	ratios := make([]float64, len(distances))
	for i, pt := range pointsCoordinates {
		distance := float64(distances[i])
		calcDistance := pt.DistanceToCoordinatesfloat32(x, y)
		ratios[i] = calcDistance / distance
	}

	// calculate median of the ratios list
	errorRatio, err = stats.Median(ratios)
	if err != nil {
		return 0, fmt.Errorf("error calculating median erro ratio. Ratios: %v.\n\tTrace: %s", ratios, err.Error())
	}
	// adjust ratio value
	errorRatio = errorRatio - 1
	return errorRatio, err
}

func HasTorableDiffByDynamicScale(a float64, b float64, toleranceDiffRatio float64) bool {

	// check if are equals
	if a == b {
		return true
	}

	// calculates diff in absolute values
	absA := math.Abs(a)
	absB := math.Abs(b)

	// if A and B has the same sign then Diff is a subtraction
	diff := math.Abs(absA - absB)

	if (a >= 0 && b < 0) || (a < 0 && b >= 0) {
		// if A and B has diferent sign then Diff is an addition
		diff = absA + absB
	}

	// use A as refence value to calcuate diff ratio
	denominator := absA

	if denominator == 0 {
		// Prevent divide by 0 (zero)
		denominator = absB
	}

	// evaluates if the Difference is bigger than A or the relation is smaller or equal to the toletared ratio
	return diff < absA && (diff/denominator) <= toleranceDiffRatio
}

//
// Calculates location by trilateration math method
//
// The ecuations used in this function are the 3 circles ecuation system
//
// r1^2 = x^2 + y^2
//
// r2^2 = (x-d)^2 + y^2
//
// r3^2 = (x-i)^2 + (y-j)^2
//
// Then, the Trilateralation coordinates calcualtion ecuation system (given by 3 circles interection)
//
// x = (r1^2 - r2^2 + d^2) / (2*d)
//
// y = ((r1^2 - r3^2 + i^2 + j^2)/2*j) - (i * x/j)
//
// This calculation includes an error correction that ensures the equations simplification of the method works,
// when doesn't have two points aligned properly the axes rotation must apply and impacts to the other
// two points that aren't the coordinate origin in the method.
//
// To calculate d,i and j variables uses traslation and rotation axes, see Point.TranslationTo Point.RotateAxesTo and Point.InvertAxesRotationTo
//
// input: the distances to the points coordinates
// output: x and y calculated location coordinates
//
// For more information please see https://en.wikipedia.org/wiki/True-range_multilateration#Three_Cartesian_dimensions.2C_three_measured_slant_ranges
func CalculateLocationByTrilateration(distances []float32, pointsCoodrinates []model.Point) (x, y float32) {
	// get distances as radius from each point converted to float64 to use math library
	radiusP1 := float64(distances[0])
	radiusP2 := float64(distances[1])
	radiusP3 := float64(distances[2])

	// get points coorinates
	p1 := pointsCoodrinates[0]
	p2 := pointsCoodrinates[1]
	p3 := pointsCoodrinates[2]

	// apply axes translation to P1 as center of cartesian axes. So circle ecuation for P1 circle keeps simple
	//p1Prime := Point{x: 0, y: 0} -- not used, just for reference
	p2Prime := p2.TranslationTo(p1)
	p3Prime := p3.TranslationTo(p1)

	// Calculate alfa angle to rotate axes. Aligning p1Prime with p2Prime points to use it as X" axis
	axesRotationAngle := math.Atan(p2Prime.Y / p2Prime.X)

	// rotates p2Prime to the axes rotation angle
	p22ndPrime := p2Prime.RotateAxesTo(axesRotationAngle)

	// rotates p3Prime to the axes rotation angle
	p32ndPrime := p3Prime.RotateAxesTo(axesRotationAngle)

	// set variables for trilateralation formula
	d := p22ndPrime.X
	i := p32ndPrime.X
	j := p32ndPrime.Y

	var resultTrilateralation model.Point

	// calculate X coordinate
	xNumerator := math.Pow(radiusP1, 2) - math.Pow(radiusP2, 2) + math.Pow(d, 2)
	xDenominator := 2 * d
	resultTrilateralation.X = xNumerator / xDenominator

	// calculate Y ecuation terms (1st and 2nd terms)
	y1stTermNumerator := math.Pow(radiusP1, 2) - math.Pow(radiusP3, 2) + math.Pow(i, 2) + math.Pow(j, 2)
	y1stTermDenominator := 2 * j

	y2ndTermNumerator := i * resultTrilateralation.X
	y2ndTermDenominator := j

	// calculate Y coordinate
	resultTrilateralation.Y = (y1stTermNumerator / y1stTermDenominator) - (y2ndTermNumerator / y2ndTermDenominator)

	traslatedLoc := resultTrilateralation.InvertAxesRotationTo(axesRotationAngle)

	// calculate x location coordinate to original reference that was traslated from P1
	x = float32(traslatedLoc.X + p1.X)

	// calculate y location coordinate to original reference that was traslated from P1
	y = float32(traslatedLoc.Y + p1.Y)

	return x, y
}
