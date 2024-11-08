package geometry

import "fmt"

type IndexedTriangle struct {
	I1, I2, I3 int
}

func MakeTriangle(I1, I2, I3 int) IndexedTriangle {
	return IndexedTriangle{
		I1,
		I2,
		I3,
	}
}

func (t IndexedTriangle) GetVectors(points []Vector3) (Vector3, Vector3, Vector3) {
	return points[t.I1], points[t.I2], points[t.I3]
}

func (t IndexedTriangle) ToFormat(points []Vector3) string {
	P1, P2, P3 := t.GetVectors(points)
	return fmt.Sprintf("[%s, %s, %s]", P1.ToFormat(), P2.ToFormat(), P3.ToFormat())
}
