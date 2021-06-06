package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	FrameThickness     = 0.15
	FrameColorEpsilon  = 0.02
	FrameWidth         = 2.3
	FrameBackDepth     = 0.6
	FrameFrontDepth    = 0.8
	FrameHeight        = 3.8
	FrameRailThickness = 0.15
)

type Frame struct {
	model3d.Solid
}

func NewFrame() Frame {
	box := &model3d.SubtractedSolid{
		Positive: model3d.NewRect(
			model3d.XYZ(-(FrameWidth+FrameThickness), -FrameFrontDepth, -FrameThickness),
			model3d.XYZ(FrameWidth+FrameThickness, (FrameBackDepth+FrameThickness),
				FrameHeight+FrameThickness),
		),
		Negative: model3d.NewRect(
			model3d.XYZ(-FrameWidth, -(FrameFrontDepth+0.01), 0.0),
			model3d.XYZ(FrameWidth, FrameBackDepth, FrameHeight),
		),
	}
	return Frame{
		Solid: model3d.JoinedSolid{
			box,
			model3d.TranslateSolid(
				hangRail(),
				model3d.YZ(FrameBackDepth+FrameThickness, FrameHeight+FrameThickness),
			),
			model3d.TranslateSolid(
				model3d.TransformSolid(
					model3d.Rotation(model3d.Y(1), math.Pi),
					hangRail(),
				),
				model3d.YZ(FrameBackDepth+FrameThickness, -FrameThickness),
			),
		},
	}
}

func (f Frame) Color(c model3d.Coord3D) render3d.Color {
	min, max := f.Min(), f.Max()
	if c.Y > FrameBackDepth+FrameThickness/2 ||
		math.Abs(model3d.NewRect(min, max).SDF(c)) < FrameColorEpsilon {
		return render3d.NewColorRGB(1.0, 0.84, 0.0)
	}
	if c.Z <= min.Z+FrameThickness+FrameColorEpsilon {
		return render3d.NewColorRGB(0.1, 0.8, 0.1)
	}
	return render3d.NewColorRGB(0.0, 190.0/255, 1.0)
}

func hangRail() model3d.Solid {
	w := FrameWidth + FrameThickness
	return model3d.CheckedFuncSolid(
		model3d.XYZ(-w, 0, -FrameRailThickness*2),
		model3d.XYZ(w, FrameRailThickness, 0),
		func(c model3d.Coord3D) bool {
			if c.Z > -FrameRailThickness {
				return true
			}
			return c.Y >= -FrameRailThickness-c.Z
		},
	)
}
