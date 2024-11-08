package screen

import (
	"fmt"
	"math"
)

const (
	BrightnessCharacter = ".:-=+*%@#"
)

type ScreenInfo struct {
	Height, Width            int
	PlaneRangeX, PlaneRangeY float64
}

type ScreenPosition struct {
	X, Y int
}

type ScreenTriangle struct {
	P1, P2, P3 ScreenPosition
	Brightness float64
}

func MakeScreenPosition(x, y int) ScreenPosition {
	return ScreenPosition{
		x,
		y,
	}
}

func MakeScreenTriangle(P1, P2, P3 ScreenPosition, Brightness float64) ScreenTriangle {
	return ScreenTriangle{
		P1,
		P2,
		P3,
		Brightness,
	}
}

func (t ScreenTriangle) sign() int {
	return (t.P1.X-t.P3.X)*(t.P2.Y-t.P3.Y) - (t.P2.X-t.P3.X)*(t.P1.Y-t.P3.Y)
}

func (t ScreenTriangle) pInTriangle(pos ScreenPosition) bool {
	var hasNeg, hasPos bool

	d1 := MakeScreenTriangle(pos, t.P1, t.P2, 0).sign()
	d2 := MakeScreenTriangle(pos, t.P2, t.P3, 0).sign()
	d3 := MakeScreenTriangle(pos, t.P3, t.P1, 0).sign()

	hasNeg = (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos = (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

func (t ScreenTriangle) getPositionsInTriangle() []ScreenPosition {
	// Find smallest and highest X and Y
	fMinX := float64(t.P1.X)
	fMinX = math.Min(fMinX, float64(t.P2.X))
	fMinX = math.Min(fMinX, float64(t.P3.X))
	minX := int(fMinX)

	fMinY := float64(t.P1.Y)
	fMinY = math.Min(fMinY, float64(t.P2.Y))
	fMinY = math.Min(fMinY, float64(t.P3.Y))
	minY := int(fMinY)

	fMaxX := float64(t.P1.X)
	fMaxX = math.Max(fMaxX, float64(t.P2.X))
	fMaxX = math.Max(fMaxX, float64(t.P3.X))
	maxX := int(fMaxX)

	fMaxY := float64(t.P1.Y)
	fMaxY = math.Max(fMaxY, float64(t.P2.Y))
	fMaxY = math.Max(fMaxY, float64(t.P3.Y))
	maxY := int(fMaxY)

	positions := make([]ScreenPosition, 0)
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			newPos := MakeScreenPosition(x, y)
			if !t.pInTriangle(newPos) {
				continue
			}

			positions = append(positions, newPos)
		}
	}

	return positions
}

func (p ScreenPosition) consoleGoToPos() {
	fmt.Printf("\033[%d;%dH", p.X, p.Y)
}

func DrawScreen(screenInfo ScreenInfo, triangles []ScreenTriangle) {
	fmt.Print("\033[H\033[2J") // clear screen
	fmt.Print("\033[?25l")     // hide cursor

	for _, triangle := range triangles {
		brightnessCharacterIdx := int(triangle.Brightness * float64(len(BrightnessCharacter)-1))
		for _, position := range triangle.getPositionsInTriangle() {
			// this position should be drawn
			if position.X < 0 || position.Y < 0 {
				continue
			}

			position.consoleGoToPos()
			fmt.Print(string(BrightnessCharacter[brightnessCharacterIdx]))
		}
	}
}
