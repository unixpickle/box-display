package main

import (
	"os"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

type Penguin struct {
	model3d.Solid
}

func NewPenguin() *Penguin {
	r, err := os.Open("models/penguin/files/Penguin_t.stl")
	essentials.Must(err)
	triangles, err := model3d.ReadSTL(r)
	essentials.Must(err)
	mesh := model3d.NewMeshTriangles(triangles).Scale(1.0 / 40.0)
	// Center the mesh on x/y, start it at z=0.
	mesh = mesh.Translate(mesh.Min().Mid(mesh.Max()).Scale(-1))
	mesh = mesh.Translate(model3d.Z(-mesh.Min().Z))
	return &Penguin{
		Solid: model3d.NewColliderSolid(model3d.MeshToCollider(mesh)),
	}
}

func (p *Penguin) Color(c model3d.Coord3D) render3d.Color {
	if c.Y < 0 {
		// Front patch of white.
		proj := c.XZ()
		center := model2d.XY(0, 0.8)
		radii := model2d.XY(0.2, 0.4)
		if center.Div(radii).Dist(proj.Div(radii)) < 1.0 {
			return render3d.NewColor(1.0)
		}
	}
	return render3d.NewColor(0.1)
}
