package main

import (
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

type Frame struct {
	model3d.Solid
}

func NewFrame() Frame {
	return Frame{
		Solid: &model3d.SubtractedSolid{
			Positive: model3d.NewRect(model3d.XYZ(-2.5, -0.8, -0.2), model3d.XYZ(2.5, 0.8, 5.0)),
			Negative: model3d.NewRect(model3d.XYZ(-2.3, -0.81, 0.0), model3d.XYZ(2.3, 0.81, 4.8)),
		},
	}
}

func (f Frame) Color(model3d.Coord3D) render3d.Color {
	return render3d.NewColorRGB(0.5, 0.5, 1.0)
}
