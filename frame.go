package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	FrameThickness    = 0.15
	FrameColorEpsilon = 0.02
	FrameWidth        = 2.3
	FrameBackDepth    = 0.6
	FrameFrontDepth   = 0.8
	FrameHeight       = 3.8
)

type Frame struct {
	model3d.Solid
}

func NewFrame() Frame {
	return Frame{
		Solid: &model3d.SubtractedSolid{
			Positive: model3d.NewRect(
				model3d.XYZ(-(FrameWidth+FrameThickness), -FrameFrontDepth, -FrameThickness),
				model3d.XYZ(FrameWidth+FrameThickness, (FrameBackDepth+FrameThickness),
					FrameHeight+FrameThickness),
			),
			Negative: model3d.NewRect(
				model3d.XYZ(-FrameWidth, -(FrameFrontDepth+0.01), 0.0),
				model3d.XYZ(FrameWidth, FrameBackDepth, FrameHeight),
			),
		},
	}
}

func (f Frame) Color(c model3d.Coord3D) render3d.Color {
	min, max := f.Min(), f.Max()
	if math.Abs(model3d.NewRect(min, max).SDF(c)) < FrameColorEpsilon {
		return render3d.NewColorRGB(1.0, 0.84, 0.0)
	}
	if c.Z <= min.Z+FrameThickness+FrameColorEpsilon {
		return render3d.NewColorRGB(0.1, 0.8, 0.1)
	}
	return render3d.NewColorRGB(0.0, 190.0/255, 1.0)
}
