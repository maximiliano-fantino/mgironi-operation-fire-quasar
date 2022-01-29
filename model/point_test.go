package model_test

import (
	"math"
	"testing"

	"github.com/mgironi/operation-fire-quasar/model"
)

// Test point translation simple
func TestPointTranslationTo(t *testing.T) {
	fatalMsg := "Traslation result didn't match. Is (%.2f.., %.2f) want (%.2f..,%.2f..)."

	type testData struct {
		point     model.Point
		newOrigin model.Point
		wanted    model.Point
	}

	testDataSet := []testData{
		{point: model.Point{X: 0, Y: 0}, newOrigin: model.Point{X: 0, Y: 0}, wanted: model.Point{X: 0, Y: 0}},
		{point: model.Point{X: 100, Y: 100}, newOrigin: model.Point{X: -50, Y: 50}, wanted: model.Point{X: 150, Y: 50}},
		{point: model.Point{X: 100, Y: 100}, newOrigin: model.Point{X: 50, Y: -50}, wanted: model.Point{X: 50, Y: 150}},
	}

	for _, tstData := range testDataSet {
		result := tstData.point.TranslationTo(tstData.newOrigin)
		if !result.EqualTo(tstData.wanted) {
			t.Errorf(fatalMsg, result.X, result.Y, tstData.wanted.X, tstData.wanted.Y)
		}
	}
}

func TestPointDistanceToPoint(t *testing.T) {
	pointA := model.Point{X: 100, Y: 0}
	pointB := model.Point{X: 200, Y: 0}
	distance := pointA.DistanceToPoint(pointB)
	wanted := float64(100)
	if distance != wanted {
		t.Errorf(`Distance calc result didn't match. Is %f want %f`, distance, wanted)
	}
}

func TestPointRotateAxesTo(t *testing.T) {
	point := model.Point{X: 100, Y: 100}
	rotatedPoint := point.RotateAxesTo(0)
	wanted := model.Point{X: 100, Y: 100}
	if !rotatedPoint.EqualTo(wanted) {
		t.Errorf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)

	}

	point = model.Point{X: 100, Y: 100}
	rotatedPoint = point.RotateAxesTo(math.Pi / 4)
	wanted = model.Point{X: math.Sqrt(math.Pow(point.X, 2) + math.Pow(point.Y, 2)), Y: 0}
	if !rotatedPoint.EqualTo(wanted) {
		t.Errorf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

}

func TestPointInvertAxesRotationTo(t *testing.T) {
	point := model.Point{X: 100, Y: 100}
	rotatedPoint := point.InvertAxesRotationTo(0)
	wanted := model.Point{X: 100, Y: 100}
	if !rotatedPoint.EqualTo(wanted) {
		t.Errorf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

	point = model.Point{X: 100, Y: 100}
	rotatedPoint = point.RotateAxesTo(math.Pi / 4)
	rotInvertedPoint := rotatedPoint.InvertAxesRotationTo(math.Pi / 4)
	wanted = point
	if !rotInvertedPoint.EqualTo(wanted) {
		t.Errorf(`Rotation calc result didn't match. Is (%f, %f) want (%f, %f)`, rotatedPoint.X, rotatedPoint.Y, wanted.X, wanted.Y)
	}

}
