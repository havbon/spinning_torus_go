package geometry

import (
	"math"

	"github.com/havbon/spinning_torus/geometry"
)

type GenerationData struct {
	Center      geometry.Vector3
	MajorRadius float64
	MinorRadius float64
}

func GetOnePoint(genData GenerationData, centralAxisAngle, tubeAngle float64) geometry.Vector3 {
	return geometry.MakeVector(
		(genData.MajorRadius+genData.MinorRadius*math.Cos(tubeAngle))*math.Cos(centralAxisAngle),
		(genData.MajorRadius+genData.MinorRadius*math.Cos(tubeAngle))*math.Sin(centralAxisAngle),
		genData.MinorRadius*math.Sin(tubeAngle),
	)
}

func GetPoints(genData GenerationData, centralResolution, tubeResolution int, zRot float64) ([]geometry.Vector3, []geometry.IndexedTriangle) {
	var points []geometry.Vector3 = make([]geometry.Vector3, 0)
	var triangles []geometry.IndexedTriangle = make([]geometry.IndexedTriangle, 0)

	cAngle := 0.0
	tAngle := 0.0
	for i := 0; i < centralResolution; i++ {
		for j := 0; j < tubeResolution; j++ {
			// generate angles for i and j
			cAngle = float64(i) / float64(centralResolution) * math.Pi * 2.0
			tAngle = float64(j) / float64(tubeResolution) * math.Pi * 2.0

			// generate another point shifted by half a step (as if i + 0.5, j + 0.5)
			cHalfStep := 0.5 / float64(centralResolution) * math.Pi * 2.0
			tHalfStep := 0.5 / float64(tubeResolution) * math.Pi * 2.0

			// find and save 4 unique points
			p1 := GetOnePoint(genData, cAngle, tAngle)
			p2 := GetOnePoint(genData, cAngle, tAngle+tHalfStep)
			p3 := GetOnePoint(genData, cAngle+cHalfStep, tAngle)
			p4 := GetOnePoint(genData, cAngle+cHalfStep, tAngle+tHalfStep)

			// apply z rotation
			p1 = p1.RotateY(zRot)
			p2 = p2.RotateY(zRot)
			p3 = p3.RotateY(zRot)
			p4 = p4.RotateY(zRot)

			// adjust to the center as per generation data
			p1 = p1.Add(genData.Center)
			p2 = p2.Add(genData.Center)
			p3 = p3.Add(genData.Center)
			p4 = p4.Add(genData.Center)

			// append points to slice
			points = append(points, p1)
			points = append(points, p2)
			points = append(points, p3)
			points = append(points, p4)

			// create triangles
			lastIdx := len(points) - 1
			tri1 := geometry.MakeTriangle(lastIdx, lastIdx-1, lastIdx-2)
			tri2 := geometry.MakeTriangle(lastIdx-1, lastIdx-3, lastIdx-2)

			// append triangles
			triangles = append(triangles, tri1)
			triangles = append(triangles, tri2)
		}
	}

	return points, triangles
}
