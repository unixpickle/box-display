package main

import (
	"math"

	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	VaseColorEpsilon = 0.02
	VaseHoleInset    = 0.22
	VaseHoleZ        = 0.9
)

type Vase struct {
	model3d.Solid
	holeRadius float64
}

func NewVase() *Vase {
	bezier := model2d.BezierCurve{
		model2d.XY(0.3, 0.0),
		model2d.XY(0.7, 0.4),
		model2d.XY(0.0, 0.7),
		model2d.XY(0.4, 1.0),
	}.Transpose()
	solid2d := model2d.CheckedFuncSolid(
		model2d.XY(0, 0),
		model2d.XY(1.0, 1.0),
		func(c model2d.Coord) bool {
			return math.Abs(c.X) < bezier.EvalX(c.Y)
		},
	)
	holeSolid2d := model2d.CheckedFuncSolid(
		model2d.XY(0, VaseHoleZ),
		model2d.XY(1.0, 1.0),
		func(c model2d.Coord) bool {
			return math.Abs(c.X) < bezier.EvalX(c.Y)-VaseHoleInset
		},
	)
	return &Vase{
		Solid: &model3d.SubtractedSolid{
			Positive: model3d.RevolveSolid(solid2d, model3d.Z(1)),
			Negative: model3d.RevolveSolid(holeSolid2d, model3d.Z(1)),
		},
		holeRadius: bezier.EvalX(VaseHoleZ) - VaseHoleInset,
	}
}

func (v *Vase) Color(c model3d.Coord3D) render3d.Color {
	if c.XY().Norm() <= v.holeRadius && math.Abs(c.Z-VaseHoleZ) < VaseColorEpsilon {
		// Dirt color.
		return render3d.NewColorRGB(0.29, 0.23, 0.0)
	} else {
		// Clay color.
		return render3d.NewColorRGB(0.74, 0.38, 0.26)
	}
}
