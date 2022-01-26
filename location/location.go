package location

import (
	"fmt"

	"log"

	"math"

	"github.com/montanaflynn/stats"
)

// Defines comparison tolerance between floats values.
const FLOAT_COMPARISION_TOLERANCE float64 = 0.0001

// Defines Point Struct with X and Y properties.
type Point struct {
	X float64
	Y float64
}

// stringfy Point value properties.
func (pt Point) String() string {
	return fmt.Sprintf("(%f, %f)", pt.X, pt.Y)
}

// Calculates the distance from Point with ohter Point.
// Uses the distance between two points formula d=sqrt((Xa-Xb)^2 + (Ya-Yb)^2).
// Returns the distance.
func (pt Point) DistanceToPoint(other Point) float64 {
	return math.Sqrt(math.Pow(pt.X-other.X, 2) + math.Pow(pt.Y-other.Y, 2))
}

// Calculates the distance from Point with a coordinate x,y given in float32.
// See also DistanceToPoint.
// input: x and Y coordiantes in float32.
// output: the distance.
func (pt Point) DistanceToCoordinatesfloat32(x, y float32) float64 {
	return pt.DistanceToPoint(Point{X: float64(x), Y: float64(y)})
}

// Calcualtes Point translation to a new referencePoint.
// input: reference point to translate to.
// output: the new Point with traslation.
func (pt Point) TranslationTo(referencePoint Point) Point {
	return Point{X: (pt.X - referencePoint.X), Y: pt.Y - referencePoint.Y}
}

// Rotates axes Point coordinates to a given axes rotation angle. Turn anticlockwise.
// The coordinates rotation forumala are x'=x*cos(angle)+y*sin(angle) and y'=y*cos(angle)-x*sin(angle).
// input: axes rotation angle in radians.
// output: the new Point with rotation.
func (pt Point) RotateAxesTo(axesRotationAngle float64) Point {
	// calculates x using coordinates rotation formula. x'=x*cos(angle)+y*sin(angle).
	x := pt.X*math.Cos(axesRotationAngle) + pt.Y*math.Sin(axesRotationAngle)

	// calculate y using coordinates rotation formula. y'=y*cos(angle)-x*sin(angle).
	y := pt.Y*math.Cos(axesRotationAngle) - pt.X*math.Sin(axesRotationAngle)

	return Point{X: x, Y: y}
}

// Rotates axes Point coordinates to a given axes rotation angle with an inverse direction to RotateAxesTo. Turn clockwise.
// The coordinates rotation inverse formula are x=x'*cos(angle)-y'*sin(angle) and y=y'*cos(angle)+x'*sin(angle).
// input: axes rotation angle in radians.
// output: the the new Point with inverse rotation.
func (pt Point) InvertAxesRotationTo(axesRotationAngle float64) Point {
	// calculates x with x=x'*cos(angle)-y'*sin(angle)
	x := pt.X*math.Cos(axesRotationAngle) - pt.Y*math.Sin(axesRotationAngle)

	// calculates y with y=y'*cos(angle)+x'*sin(angle)
	y := pt.Y*math.Cos(axesRotationAngle) + pt.X*math.Sin(axesRotationAngle)

	return Point{X: x, Y: y}
}

// Checks if Point is equal with other Point, compares property by property using a float comparision tolerance.
// See also FLOAT_COMPARISION_TOLERANCE.
// input: other point
// output: true if equal otherwise false
func (pt Point) EqualTo(otherPoint Point) bool {
	diffX := math.Abs(pt.X - otherPoint.X)
	isEqualX := diffX < FLOAT_COMPARISION_TOLERANCE

	diffY := math.Abs(pt.Y - otherPoint.Y)
	isEqualY := diffY < FLOAT_COMPARISION_TOLERANCE

	return isEqualX && isEqualY
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
	acceptedRatio := FLOAT_COMPARISION_TOLERANCE
	if ratioErr > acceptedRatio {
		log.Printf("WARN ratio error exceeds aceptable level of %.4f. Ratio ~ %.4f", acceptedRatio, ratioErr)
	}
	return x, y, nil
}

// Routput: the kwnown reference coordinates.
func GetKnownReferenceCoordinates() []Point {
	return []Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}
}

// Checks if the X, Y coordinates distance to each pointsCoordinates matchs with the given distances.
// input: distances, points coordinates and 'x','y' calculated coordinates to check.
// output: the median errorRatio calculated (0: no error, interval [0,1]: percent error)
// error1: if detects arrays length diferences (between distances and points coordinates)
// error2: an internal calculation error.
func ChecksDistancesToCoordinate(distances []float32, pointsCoordinates []Point, x, y float32) (errorRatio float64, err error) {
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
		return 0, err
	}
	// adjust ratio value
	errorRatio = errorRatio - 1
	return errorRatio, err
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
func CalculateLocationByTrilateration(distances []float32, pointsCoodrinates []Point) (x, y float32) {
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

	var resultTrilateralation Point

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
