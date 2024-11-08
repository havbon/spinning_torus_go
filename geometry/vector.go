package geometry

import (
	"errors"
	"fmt"
	"math"

	"github.com/havbon/spinning_torus/screen"
)

type Vector3 struct {
	X, Y, Z float64
}

func MakeVector(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func Zero() Vector3 {
	return MakeVector(0, 0, 0)
}

func One() Vector3 {
	return MakeVector(1, 1, 1)
}

func (v Vector3) MultiplyF(val float64) Vector3 {
	return MakeVector(
		v.X*val,
		v.Y*val,
		v.Z*val,
	)
}

func (v Vector3) Add(other Vector3) Vector3 {
	return MakeVector(
		v.X+other.X,
		v.Y+other.Y,
		v.Z+other.Z,
	)
}

func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (a Vector3) Cross(b Vector3) Vector3 {
	return MakeVector(
		a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X,
	)
}

func (v Vector3) Subtract(other Vector3) Vector3 {
	return MakeVector(
		v.X-other.X,
		v.Y-other.Y,
		v.Z-other.Z,
	)
}

func (v Vector3) Length() float64 {
	return math.Pow(
		math.Pow(v.X, 2)+math.Pow(v.Y, 2)+math.Pow(v.Z, 2),
		0.5,
	)
}

func (v Vector3) Normalize() Vector3 {
	return v.MultiplyF(1 / v.Length())
}

func (v Vector3) ToFormat() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}

func (v Vector3) AngleBetween(other Vector3) float64 {
	dotProduct := v.Dot(other)
	productLength := v.Length() * other.Length()
	return math.Acos(dotProduct / productLength)
}

func (v Vector3) RotateY(angle float64) Vector3 {
	x := MakeVector(math.Cos(angle), 0, math.Sin(angle))
	y := MakeVector(0, 1, 0)
	z := MakeVector(-math.Sin(angle), 0, math.Cos(angle))

	return MakeVector(
		x.Dot(v),
		y.Dot(v),
		z.Dot(v),
	)
}

func (v Vector3) RotateX(angle float64) Vector3 {
	x := MakeVector(0, 0, 0)
	y := MakeVector(0, math.Cos(angle), -math.Sin(angle))
	z := MakeVector(0, math.Sin(angle), math.Cos(angle))

	return MakeVector(
		x.Dot(v),
		y.Dot(v),
		z.Dot(v),
	)
}

func (v Vector3) RotateZ(angle float64) Vector3 {
	x := MakeVector(math.Cos(angle), -math.Sin(angle), 0)
	y := MakeVector(math.Sin(angle), math.Cos(angle), 0)
	z := MakeVector(0, 0, 1)

	return MakeVector(
		x.Dot(v),
		y.Dot(v),
		z.Dot(v),
	)
}

func (v Vector3) PlaneIntersect(pointOnLine Vector3, planePoint Vector3, planeNormal Vector3) (Vector3, error) {
	var err error

	unitV := v.Normalize()

	denominator := unitV.Dot(planeNormal)

	if denominator == 0 {
		err = errors.New("plane and vector are parallell")
		return Zero(), err
	}

	numerator := planePoint.Subtract(pointOnLine).Dot(planeNormal)

	d := numerator / denominator

	return pointOnLine.Add(unitV.MultiplyF(d)), nil
}

func (N Vector3) ChooseNonParallelVector() Vector3 {
	const threshold = 0.99

	if math.Abs(N.X) < threshold {
		return MakeVector(1, 0, 0)
	} else if math.Abs(N.Y) < threshold {
		return MakeVector(0, 1, 0)
	} else {
		return MakeVector(0, 0, 1)
	}
}

func (projectedV Vector3) ProjectedPointTo2D(planePoint Vector3, planeNormal Vector3) Vector3 {
	notParallel := planeNormal.ChooseNonParallelVector()
	u := planeNormal.Cross(notParallel).Normalize()
	v := planeNormal.Cross(u)

	newX := projectedV.Subtract(planePoint).Dot(u)
	newY := projectedV.Subtract(planePoint).Dot(v)

	return MakeVector(
		newX,
		newY,
		0,
	)
}

func (v Vector3) MapProjection(screenInfo screen.ScreenInfo) screen.ScreenPosition {
	// shift projection, every value is positive
	shifted := v.Add(MakeVector(
		screenInfo.PlaneRangeX/2,
		screenInfo.PlaneRangeY/2,
		0,
	))
	// screen is now starting at 0 and going to plane range x/y

	// divide x and y by the total length of the screen
	// which puts every value between 0-1
	norm := MakeVector(
		shifted.X/screenInfo.PlaneRangeX,
		shifted.Y/screenInfo.PlaneRangeY,
		0,
	)

	// multiply by the screen width and height
	mapped := MakeVector(
		norm.X*float64(screenInfo.Height),
		norm.Y*float64(screenInfo.Width),
		0,
	)

	// round, int, return
	return screen.MakeScreenPosition(
		int(math.Round(mapped.X)),
		int(math.Round(mapped.Y)),
	)
}
