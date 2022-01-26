package location_test

import (
	"fmt"
	"testing"

	"math"

	"github.com/mgironi/operation-fire-quasar/location"
)

// Test 'location.CalculateLocation' with all nil distances
// Error expected as 'Is not possible to compelete calculations'
//TODO this test is for main function

// Test 'location.CalculateLocation' with some nil distances
// Error expected as 'Is not possible to compelete calculations'
//TODO this test is for main function

// Test 'location.CalculateLocation' with a simple call of real distances
func TestCalculateLocationWithSimple(t *testing.T) {
	pointsCoordinates := []location.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}
	wanted := location.Point{X: -200, Y: 200}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)
}

func TestCalculateLocationMultipleTimes(t *testing.T) {
	pointsCoordinates := []location.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}

	wanted := location.Point{X: -200, Y: 200}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: 300, Y: -700}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: -1000, Y: -900}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: 100, Y: 200}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

}

// Test 'location.CalculateLocation' with a simple call of real distances. Trimed at 2 digits presicion
func TestCalculateLocationWithSimplePrecisionTrim(t *testing.T) {
	// Trim presicion upto 2 digits from {500, 424.2640687, 707.1067812}
	distances := []float32{500, 424.26, 707.10}
	pointsCoordinates := []location.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}
	wanted := location.Point{X: -200, Y: 200}
	runAndCheckCalculateLocation(pointsCoordinates, distances, wanted, false, t)
}

// Test 'location.CalculateLocation' with a simple call of real distances. Trimed at 0 digits presicion
func TestCalculateLocationWithSimpleWOPrecision(t *testing.T) {
	// Trim presicion upto 2 digits from {500, 424.2640687, 707.1067812}
	distances := []float32{500, 424, 707}
	pointsCoordinates := []location.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}
	wanted := location.Point{X: -200, Y: 200}
	runAndCheckCalculateLocation(pointsCoordinates, distances, wanted, true, t)
}

// Test 'location.CalculateLocation' with extreme values distances.
func TestCalculateLocationMaxLimitRange(t *testing.T) {
	pointsCoordinates := []location.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}

	wanted := location.Point{X: 1000, Y: 1000}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: 10000, Y: 10000}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: 100000, Y: 100000}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, false, t)

	wanted = location.Point{X: 1000000, Y: 1000000}
	runAndCheckCalculateLocation(pointsCoordinates, nil, wanted, true, t)
}

func runAndCheckCalculateLocation(pointsCoordinates []location.Point, distances []float32, wanted location.Point, spectedError bool, t *testing.T) {
	fmt.Printf("Location wanted: (%f, %f)\n", wanted.X, wanted.Y)

	if distances == nil {
		distances = make([]float32, len(pointsCoordinates))
		for i, pt := range pointsCoordinates {
			distances[i] = float32(pt.DistanceToPoint(wanted))
			fmt.Printf("Calculated Distance to p%d (%f, %f): %f\n", i+1, pt.X, pt.Y, distances[i])
		}
	}

	var errorRatioTolerance float64 = location.FLOAT_COMPARISION_TOLERANCE

	x, y, err := location.CalculateLocation(distances)
	if err != nil {
		t.Errorf("error on calculate location. %e", err)
	}

	fmt.Printf("Calculated location (%f, %f)\n", x, y)
	errorDetcted := x != float32(wanted.X) && !hasTorableDiffByDynamicScale(float64(wanted.X), float64(x), errorRatioTolerance)

	if errorDetcted {
		if !spectedError {
			t.Fatalf(`The X coordinate don't match as spected. Is %f, want %f +/-%f%% `, x, wanted.X, 100*errorRatioTolerance)
		}
	}

	errorDetcted = y != float32(wanted.Y) && !hasTorableDiffByDynamicScale(float64(wanted.Y), float64(y), errorRatioTolerance)
	if errorDetcted {
		if !spectedError {
			t.Fatalf(`The Y coordinate don't match as spected. Is %f, want %f +/- %f%% `, y, wanted.Y, 100*errorRatioTolerance)
		}
	}

	if !errorDetcted && spectedError {
		t.Fatalf("Error was spected. Is (%f, %f), want (%f, %f) +/- %f%% ", x, y, wanted.X, wanted.Y, 100*errorRatioTolerance)
	}
}

//TODO test validation formula
func hasTorableDiffByDynamicScale(a float64, b float64, toleranceDiffRatio float64) bool {

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

func runAndCheckHasTorableDiference(a float64, b float64, toleranceDiffRatio float64, want bool, t *testing.T) {
	result := hasTorableDiffByDynamicScale(a, b, toleranceDiffRatio)
	match := result == want
	if !match {
		t.Fatalf(`The diference isn't tolerable. Between reference %f to %f with %f tolerance. Is %t, want %t.`, a, b, toleranceDiffRatio, result, want)
	}
}

func TestHasTorableDiference(t *testing.T) {
	runAndCheckHasTorableDiference(float64(0), float64(0), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(100), float64(120), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(100), float64(110), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(100), float64(105), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(-100), float64(-105), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(-100), float64(-95), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(100), float64(95), float64(0.1), true, t)
	runAndCheckHasTorableDiference(float64(5), float64(-3), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(3), float64(-5), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(0), float64(-5), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(0), float64(5), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(5), float64(0), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(-5), float64(0), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(-15), float64(10), float64(0.1), false, t)
	runAndCheckHasTorableDiference(float64(10), float64(-15), float64(0.1), false, t)
}

func TestChecksDistancesToCoordinate(t *testing.T) {
	var x, y float32 = -200, 200
	pointsCoordinates := []location.Point{
		{X: -500, Y: -200},
		{X: 100, Y: -100},
		{X: 500, Y: 100},
	}
	distances := make([]float32, 3)
	for i, pt := range pointsCoordinates {
		distances[i] = float32(pt.DistanceToCoordinatesfloat32(x, y))
	}

	errRatio, err := location.ChecksDistancesToCoordinate(distances, pointsCoordinates, x, y)
	spected := float64(0)
	if err != nil {
		t.Error(err)
	}
	match := errRatio == spected
	if !match {
		t.Fatalf(`Calculation error didn't match. Is %.2f, want %.2f.`, errRatio, spected)
	}
}

// Test point translation simple
func TestPointTranslationTo(t *testing.T) {
	fatalMsg := "Traslation result didn't match. Is (%.2f.., %.2f) want (%.2f..,%.2f..)."

	type testData struct {
		point     location.Point
		newOrigin location.Point
		wanted    location.Point
	}

	testDataSet := []testData{
		{point: location.Point{X: 0, Y: 0}, newOrigin: location.Point{X: 0, Y: 0}, wanted: location.Point{X: 0, Y: 0}},
		{point: location.Point{X: 100, Y: 100}, newOrigin: location.Point{X: -50, Y: 50}, wanted: location.Point{X: 150, Y: 50}},
		{point: location.Point{X: 100, Y: 100}, newOrigin: location.Point{X: 50, Y: -50}, wanted: location.Point{X: 50, Y: 150}},
	}

	for _, tstData := range testDataSet {
		result := tstData.point.TranslationTo(tstData.newOrigin)
		if !result.EqualTo(tstData.wanted) {
			t.Fatalf(fatalMsg, result.X, result.Y, tstData.wanted.X, tstData.wanted.Y)
		}
	}
}

func TestPointDistanceToPoint(t *testing.T) {
	pointA := location.Point{X: 100, Y: 0}
	pointB := location.Point{X: 200, Y: 0}
	distance := pointA.DistanceToPoint(pointB)
	wanted := float64(100)
	if distance != wanted {
		t.Fatalf(`Distance calc result didn't match. Is %f want %f`, distance, wanted)
	}
}

func TestRotateAxesTo(t *testing.T) {
	point := location.Point{X: 100, Y: 100}
	rotatedPoint := point.RotateAxesTo(0)
	wanted := location.Point{X: 100, Y: 100}
	if !rotatedPoint.EqualTo(wanted) {
		t.Fatalf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)

	}

	point = location.Point{X: 100, Y: 100}
	rotatedPoint = point.RotateAxesTo(math.Pi / 4)
	wanted = location.Point{X: math.Sqrt(math.Pow(point.X, 2) + math.Pow(point.Y, 2)), Y: 0}
	if !rotatedPoint.EqualTo(wanted) {
		t.Fatalf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

}

func TestInvertAxesRotationTo(t *testing.T) {
	point := location.Point{X: 100, Y: 100}
	rotatedPoint := point.InvertAxesRotationTo(0)
	wanted := location.Point{X: 100, Y: 100}
	if !rotatedPoint.EqualTo(wanted) {
		t.Fatalf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

	point = location.Point{X: 100, Y: 100}
	rotatedPoint = point.RotateAxesTo(math.Pi / 4)
	rotInvertedPoint := rotatedPoint.InvertAxesRotationTo(math.Pi / 4)
	wanted = point
	if !rotInvertedPoint.EqualTo(wanted) {
		t.Fatalf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

}
