package projection

import (
	"log"
	"math"

	"github.com/havbon/spinning_torus/geometry"
	"github.com/havbon/spinning_torus/screen"
)

type Camera struct {
	FocalPoint        geometry.Vector3
	ScreenPlaneOffset geometry.Vector3
}

type LightingSource struct {
	Position  geometry.Vector3
	Intensity float64
}

type ScreenVector struct {
	X, Y int
}

type ProjectedTriangle struct {
	P1, P2, P3 ScreenVector
}

func MakeScreenVector(x, y int) ScreenVector {
	return ScreenVector{
		x,
		y,
	}
}

func MakeProjectedTriangle(p1, p2, p3 ScreenVector) ProjectedTriangle {
	return ProjectedTriangle{
		p1,
		p2,
		p3,
	}
}

func (lightingSource LightingSource) GetTriangleBrightness(p1, p2, p3 geometry.Vector3) float64 {
	// find triangle normal
	triangleV1 := p2.Subtract(p1).Normalize()
	triangleV2 := p3.Subtract(p1).Normalize()
	triangleNormal := triangleV1.Cross(triangleV2).Normalize()

	// direction from triangle to light source normalized
	lightingDir := lightingSource.Position.Subtract(p1).Normalize()
	// angle between light and triangle normal
	cosAngle := lightingDir.Dot(triangleNormal)
	// find brightness
	brightness := lightingSource.Intensity * math.Max(0, cosAngle)
	brightness = math.Min(1, brightness)
	return brightness
}

func Project(
	screenInfo screen.ScreenInfo,
	camera Camera,
	lightingSource LightingSource,
	points []geometry.Vector3,
	indexes []geometry.IndexedTriangle,
) []screen.ScreenTriangle {
	// draw line from every triangle to focal point
	// screen plane offset works as a normal line to the plane
	// define this plane
	planePoint := camera.FocalPoint.Add(camera.ScreenPlaneOffset)
	planeNormal := camera.FocalPoint.Subtract(camera.ScreenPlaneOffset).Normalize()

	projected := make([]screen.ScreenTriangle, 0)

	for _, triangle := range indexes {
		p1, p2, p3 := points[triangle.I1], points[triangle.I2], points[triangle.I3]

		// if any of p1,p2,p3 is on the same side of plane as focal, ignore triangle
		planeToFocal := camera.FocalPoint.Subtract(planePoint).Normalize()
		planeToP1 := p1.Subtract(planePoint).Normalize()
		planeToP2 := p2.Subtract(planePoint).Normalize()
		planeToP3 := p3.Subtract(planePoint).Normalize()

		p1Outside := math.Pi/2-planeToP1.AngleBetween(planeToFocal) < 0
		p2Outside := math.Pi/2-planeToP2.AngleBetween(planeToFocal) < 0
		p3Outside := math.Pi/2-planeToP3.AngleBetween(planeToFocal) < 0

		if p1Outside || p2Outside || p3Outside {
			continue
		}

		focalToP1 := p1.Subtract(camera.FocalPoint)
		focalToP2 := p2.Subtract(camera.FocalPoint)
		focalToP3 := p3.Subtract(camera.FocalPoint)

		// these are the vectors scaled so they hit the plane
		planeIntersect1, err1 := focalToP1.PlaneIntersect(camera.FocalPoint, planePoint, planeNormal)
		planeIntersect2, err2 := focalToP2.PlaneIntersect(camera.FocalPoint, planePoint, planeNormal)
		planeIntersect3, err3 := focalToP3.PlaneIntersect(camera.FocalPoint, planePoint, planeNormal)

		if err1 != nil {
			log.Fatal(err1.Error())
		}
		if err2 != nil {
			log.Fatal(err2.Error())
		}
		if err3 != nil {
			log.Fatal(err3.Error())
		}

		// turn projection into 2D coords (still in 3D vector)
		// z = 0
		projection1 := planeIntersect1.ProjectedPointTo2D(planePoint, planeNormal)
		projection2 := planeIntersect2.ProjectedPointTo2D(planePoint, planeNormal)
		projection3 := planeIntersect3.ProjectedPointTo2D(planePoint, planeNormal)

		// map projections onto screen
		screenPos1 := projection1.MapProjection(screenInfo)
		screenPos2 := projection2.MapProjection(screenInfo)
		screenPos3 := projection3.MapProjection(screenInfo)

		brightness := lightingSource.GetTriangleBrightness(p1, p2, p3)
		projected = append(projected, screen.MakeScreenTriangle(
			screenPos1,
			screenPos2,
			screenPos3,
			brightness,
		))
	}

	return projected
}
