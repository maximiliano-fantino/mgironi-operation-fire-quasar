package model

import (
	"fmt"
	"math"
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
