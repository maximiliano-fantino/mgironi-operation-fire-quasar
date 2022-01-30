package location_test

import (
	"fmt"
	"testing"

	test "github.com/mgironi/operation-fire-quasar/_test"
	"github.com/mgironi/operation-fire-quasar/location"
	"github.com/mgironi/operation-fire-quasar/model"
)

// Test 'location.CalculateLocation' with all nil distances
// Error expected as 'Is not possible to compelete calculations'
func TestGetLocationEmptyDistances(t *testing.T) {
	distances := []float32{}
	wantedX := float32(0)
	wantedY := float32(0)
	gotX, gotY := location.GetLocation(distances...)

	if !test.AreFloats32Equals(gotX, wantedX) {
		t.Errorf("Error GetLocation() for X coordinate, got: %f wanted:%f", gotX, wantedY)
	}

	if !test.AreFloats32Equals(gotY, wantedY) {
		t.Errorf("Error GetLocation() for Y coordinate, got: %f wanted:%f", gotY, wantedY)
	}
}

// Test 'location.CalculateLocation' with some nil distances
// Error expected as 'Is not possible to compelete calculations'
func TestGetLocationSomeDistances(t *testing.T) {
	distances := []float32{100, 200}
	wantedX := float32(0)
	wantedY := float32(0)
	gotX, gotY := location.GetLocation(distances...)

	if !test.AreFloats32Equals(gotX, wantedX) {
		t.Errorf("Error GetLocation() for X coordinate, got: %f wanted:%f", gotX, wantedY)
	}

	if !test.AreFloats32Equals(gotY, wantedY) {
		t.Errorf("Error GetLocation() for Y coordinate, got: %f wanted:%f", gotY, wantedY)
	}
}

// Tests GetLocation
func TestGetLocation(t *testing.T) {
	distances := []float32{100, 200, 400}
	//	wantedX := float32(0)
	//	wantedY := float32(0)
	location.GetLocation(distances...)

	/*	if test.AreFloats32Equals(gotX, wantedX) {
			t.Errorf("Error GetLocation() for X coordinate, got: %f wanted:%f", gotX, wantedY)
		}

		if test.AreFloats32Equals(gotY, wantedY) {
			t.Errorf("Error GetLocation() for Y coordinate, got: %f wanted:%f", gotY, wantedY)
		}
	*/
}

// Tests checkDistanceToCoordinate
func TestChecksDistancesToCoordinate(t *testing.T) {
	var x, y float32 = -200, 200
	pointsCoordinates := []model.Point{
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
		t.Errorf(`Calculation error didn't match. Is %.2f, want %.2f.`, errRatio, spected)
	}
}

// Tests HasTorableDiffByDynamicScale()
func TestHasTorableDiffByDynamicScale(t *testing.T) {
	type args struct {
		a                  float64
		b                  float64
		toleranceDiffRatio float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test1", args: args{a: 0, b: 0, toleranceDiffRatio: 0.1}, want: true},
		{name: "test2", args: args{a: 100, b: 120, toleranceDiffRatio: 0.1}, want: false},
		{name: "test3", args: args{a: 100, b: 110, toleranceDiffRatio: 0.1}, want: true},
		{name: "test4", args: args{a: 100, b: 105, toleranceDiffRatio: 0.1}, want: true},
		{name: "test5", args: args{a: -100, b: -105, toleranceDiffRatio: 0.1}, want: true},
		{name: "test6", args: args{a: -100, b: -95, toleranceDiffRatio: 0.1}, want: true},
		{name: "test7", args: args{a: 100, b: 95, toleranceDiffRatio: 0.1}, want: true},
		{name: "test8", args: args{a: 5, b: -3, toleranceDiffRatio: 0.1}, want: false},
		{name: "test9", args: args{a: 3, b: -5, toleranceDiffRatio: 0.1}, want: false},
		{name: "test10", args: args{a: 0, b: -5, toleranceDiffRatio: 0.1}, want: false},
		{name: "test11", args: args{a: 0, b: 5, toleranceDiffRatio: 0.1}, want: false},
		{name: "test12", args: args{a: 5, b: 0, toleranceDiffRatio: 0.1}, want: false},
		{name: "test13", args: args{a: -5, b: 0, toleranceDiffRatio: 0.1}, want: false},
		{name: "test14", args: args{a: -15, b: 10, toleranceDiffRatio: 0.1}, want: false},
		{name: "test15", args: args{a: 10, b: -15, toleranceDiffRatio: 0.1}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := location.HasTorableDiffByDynamicScale(tt.args.a, tt.args.b, tt.args.toleranceDiffRatio); got != tt.want {
				t.Errorf(`Test '%s' Is %t, want %t. The diference isn't tolerable. Between reference %f to %f with %f tolerance.`, tt.name, got, tt.want, tt.args.a, tt.args.b, tt.args.toleranceDiffRatio)
			}
		})
	}
}

// Tests CalculateLocation
func TestCalculateLocation(t *testing.T) {
	// cleans variables envs to prevent use of not controlled default satelite info
	test.CleanSatelitesInfoEnvs()

	// define own points coordinates
	pointsCoordinates := []model.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}

	type args struct {
		distances []float32
	}
	tests := []struct {
		name    string
		args    args
		wantX   float32
		wantY   float32
		wantErr bool
	}{
		{name: "test1", args: args{distances: nil}, wantX: -200, wantY: 200, wantErr: false},
		{name: "test2", args: args{distances: nil}, wantX: 1000, wantY: 1000, wantErr: false},
		{name: "test3", args: args{distances: nil}, wantX: -200, wantY: 200, wantErr: false},
		{name: "test4", args: args{distances: nil}, wantX: 300, wantY: -700, wantErr: false},
		{name: "test5", args: args{distances: nil}, wantX: -1000, wantY: -900, wantErr: false},
		{name: "test6", args: args{distances: nil}, wantX: 100, wantY: 200, wantErr: false},

		// Trim presicion upto 2 digits from {500, 424.2640687, 707.1067812}
		{name: "test7", args: args{distances: []float32{500, 424.26, 707.10}}, wantX: -200, wantY: 200, wantErr: false},

		// triming all precision
		{name: "test8", args: args{distances: []float32{500, 424, 707}}, wantX: -200, wantY: 200, wantErr: true},

		// out of range limit
		{name: "test9", args: args{distances: nil}, wantX: 1000, wantY: 1000, wantErr: false},
		{name: "test10", args: args{distances: nil}, wantX: 10000, wantY: 10000, wantErr: false},
		{name: "test11", args: args{distances: nil}, wantX: 100000, wantY: 100000, wantErr: false},
		{name: "test12", args: args{distances: nil}, wantX: 1000000, wantY: 1000000, wantErr: false},
	}

	// defines error ratio tolerance for check calculation
	errorRatioTolerance := model.FLOAT_COMPARISION_TOLERANCE

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgCtxErr := fmt.Sprintf("Context for test '%s' info:\n\tLocation wanted: (%f, %f)\n", tt.name, tt.wantX, tt.wantY)
			// autocomplete distances if came on nil based on wanted X and Y and poitnsCoordinates
			if tt.args.distances == nil {
				tt.args.distances = make([]float32, len(pointsCoordinates))
				for i, pt := range pointsCoordinates {
					tt.args.distances[i] = float32(pt.DistanceToPoint(model.Point{X: float64(tt.wantX), Y: float64(tt.wantY)}))
					msgCtxErr += fmt.Sprintf("\tCalculated Distance to p%d (%f, %f): %f\n", i+1, pt.X, pt.Y, tt.args.distances[i])
				}
			}
			gotX, gotY, err := location.CalculateLocation(tt.args.distances)
			if (err != nil) && !tt.wantErr {
				t.Errorf("CalculateLocation() error = %v, wantErr %v.\n%s", err, tt.wantErr, msgCtxErr)
				return
			}
			msgCtxErr += fmt.Sprintf("\tCalculated location (%f, %f)\n", gotX, gotY)

			errorDetectedX := test.AreFloats32Equals(gotX, tt.wantX) && !location.HasTorableDiffByDynamicScale(float64(tt.wantX), float64(gotX), errorRatioTolerance)
			if errorDetectedX && !tt.wantErr {
				t.Errorf(`\tCalculateLocation() gotX = %f, want %f +/-%f%%\n%s`, gotX, tt.wantX, 100*errorRatioTolerance, msgCtxErr)
			}

			errorDetectedY := test.AreFloats32Equals(gotY, tt.wantY) && !location.HasTorableDiffByDynamicScale(float64(tt.wantY), float64(gotY), errorRatioTolerance)
			if errorDetectedY && !tt.wantErr {
				t.Errorf(`\tCalculateLocation() gotY = %f, want %f +/-%f%%\n%s`, gotY, tt.wantY, 100*errorRatioTolerance, msgCtxErr)
			}

			if !(err == nil || errorDetectedX || errorDetectedY) && tt.wantErr {
				t.Errorf("\tError was spected. Is (%f, %f), want (%f, %f) +/- %f%%\n%s", gotX, gotY, tt.wantX, tt.wantY, 100*errorRatioTolerance, msgCtxErr)
			}
		})
	}
}
