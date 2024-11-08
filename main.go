package main

import (
	"github.com/eiannone/keyboard"
	"github.com/havbon/spinning_torus/geometry"
	torus "github.com/havbon/spinning_torus/geometry/torus"
	"github.com/havbon/spinning_torus/input"
	"github.com/havbon/spinning_torus/projection"
	"github.com/havbon/spinning_torus/screen"
)

func main() {
	torusGeneration := torus.GenerationData{
		Center:      geometry.Zero(),
		MajorRadius: 3,
		MinorRadius: 1,
	}
	camera := projection.Camera{
		FocalPoint:        geometry.MakeVector(0, 0, 25),
		ScreenPlaneOffset: geometry.MakeVector(0, 0, 3),
	}
	lightingSource := projection.LightingSource{
		Position:  geometry.MakeVector(0, -15, 15),
		Intensity: 0.9,
	}
	viewport := screen.ScreenInfo{
		Height:      21,
		Width:       150,
		PlaneRangeX: 1.0,
		PlaneRangeY: 3.0,
	}

	keyCh := make(chan keyboard.KeyEvent)
	xRot := 0.0
	running := true
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()
	go input.GetKey(keyCh)
	for running {
		// input handling
		select {
		case keyEv := <-keyCh:
			switch keyEv.Key {
			case keyboard.KeyEsc:
				running = false
			case keyboard.KeyArrowLeft:
				camera.ScreenPlaneOffset = camera.ScreenPlaneOffset.RotateY(0.04)
			case keyboard.KeyArrowRight:
				camera.ScreenPlaneOffset = camera.ScreenPlaneOffset.RotateY(-0.04)
			case keyboard.KeyArrowUp:
				camera.ScreenPlaneOffset = camera.ScreenPlaneOffset.RotateX(0.04)
			case keyboard.KeyArrowDown:
				camera.ScreenPlaneOffset = camera.ScreenPlaneOffset.RotateX(-0.04)
			}

			if keyEv.Rune == 'w' {
				camera.FocalPoint = camera.FocalPoint.Add(camera.ScreenPlaneOffset.MultiplyF(-0.3))
			}
			if keyEv.Rune == 's' {
				camera.FocalPoint = camera.FocalPoint.Add(camera.ScreenPlaneOffset.MultiplyF(0.3))
			}
			side := geometry.MakeVector(0, 1, 0).Cross(camera.ScreenPlaneOffset)
			if keyEv.Rune == 'a' {
				camera.FocalPoint = camera.FocalPoint.Add(side.MultiplyF(-0.3))
			}
			if keyEv.Rune == 'd' {
				camera.FocalPoint = camera.FocalPoint.Add(side.MultiplyF(0.3))
			}
		default:

		}

		points, triangles := torus.GetPoints(torusGeneration, 500, 50, float64(xRot))
		projected := projection.Project(
			viewport,
			camera,
			lightingSource,
			points,
			triangles,
		)
		screen.DrawScreen(viewport, projected)

		xRot += 0.05
	}
}
