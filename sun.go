package main

import (
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

type Sun struct {
	model3d.Solid
}

func NewSun() *Sun {
	return &Sun{
		model3d.IntersectedSolid{
			model3d.NewRect(
				model3d.XYZ(-0.5, -0.5, -0.5),
				model3d.XYZ(0.5, 0.0, 0.5),
			),
			&model3d.Sphere{Center: model3d.XYZ(0, 0.5, 0), Radius: 0.7},
		},
	}
}

func (s *Sun) Color(c model3d.Coord3D) render3d.Color {
	return render3d.NewColorRGB(0.8, 0.8, 0.0)
}
