package main

import (
	"math"

	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
	"github.com/unixpickle/model3d/toolbox3d"
)

type Cloud struct {
	model3d.Solid
}

func NewCloud(imageName string) *Cloud {
	bitmap := model2d.MustReadBitmap(imageName, nil)
	mesh2d := bitmap.Mesh().SmoothSq(30).Scale(1.0 / 900.0)
	mesh2d = mesh2d.Translate(mesh2d.Min().Mid(mesh2d.Max()).Scale(-1))

	sdf2d := model2d.MeshToSDF(mesh2d)
	hm := toolbox3d.NewHeightMap(sdf2d.Min(), sdf2d.Max(), 5000)
	if Production {
		hm.AddSpheresSDF(sdf2d, 2000, 0.04, 0.13)
	} else {
		hm.AddSpheresSDF(sdf2d, 500, 0.04, 0.13)
	}
	fullSolid := toolbox3d.HeightMapToSolid(hm)

	// Cut off bottom to make it flatter.
	fullSolid = model3d.TranslateSolid(fullSolid, model3d.Z(-0.05))
	solid := model3d.IntersectedSolid{
		fullSolid,
		model3d.NewRect(fullSolid.Min().Mul(model3d.XY(1, 1)), fullSolid.Max()),
	}

	return &Cloud{
		Solid: model3d.TransformSolid(model3d.Rotation(model3d.X(1), math.Pi/2), solid),
	}
}

func (c *Cloud) Color(_ model3d.Coord3D) render3d.Color {
	return render3d.NewColor(1.0)
}
