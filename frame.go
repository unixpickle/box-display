package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	FrameThickness    = 0.2
	FrameColorEpsilon = 0.02
)

type Frame struct {
	model3d.Solid
}

func NewFrame() Frame {
	return Frame{
		Solid: &model3d.SubtractedSolid{
			Positive: model3d.NewRect(
				model3d.XYZ(-2.5, -0.8, -FrameThickness),
				model3d.XYZ(2.5, 0.8, 4.0),
			),
			Negative: model3d.NewRect(
				model3d.XYZ(-2.5+FrameThickness, -0.81, 0.0),
				model3d.XYZ(2.5-FrameThickness, 0.8-FrameThickness, 4.0-FrameThickness),
			),
		},
	}
}

func (f Frame) Color(c model3d.Coord3D) render3d.Color {
	min, max := f.Min(), f.Max()
	if math.Abs(model3d.NewRect(min, max).SDF(c)) < FrameColorEpsilon {
		return render3d.NewColor(0.5)
	}
	if c.Z <= min.Z+FrameThickness+FrameColorEpsilon {
		return render3d.NewColorRGB(0.1, 0.8, 0.1)
	}
	return render3d.NewColorRGB(0.0, 190.0/255, 1.0)
}
