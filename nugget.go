package main

import (
	"math"

	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	NuggetRounding    = 0.1
	NuggetCrumbRadius = 0.02
)

type Nugget struct {
	mainBody model3d.SDF
	centers  *model3d.CoordTree
	colors   map[model3d.Coord3D]render3d.Color
}

func NewNugget() *Nugget {
	outline := model2d.MeshToSDF(nuggetOutline())

	mainBody := model3d.ProfileSDF(outline, -0.1, 0.1)
	mainBody = model3d.TransformSDF(model3d.Rotation(model3d.X(1), math.Pi/2), mainBody)

	// TODO: create many small spheres around the nugget
	// to simulate bread crumbs.

	return &Nugget{
		mainBody: mainBody,
	}
}

func (n *Nugget) Min() model3d.Coord3D {
	return n.mainBody.Min().Sub(model3d.XYZ(1, 1, 1).Scale(NuggetCrumbRadius + NuggetRounding))
}

func (n *Nugget) Max() model3d.Coord3D {
	return n.mainBody.Max().Add(model3d.XYZ(1, 1, 1).Scale(NuggetCrumbRadius + NuggetRounding))
}

func (n *Nugget) Contains(c model3d.Coord3D) bool {
	return n.mainBody.SDF(c) > -NuggetRounding
}

func (n *Nugget) Color(c model3d.Coord3D) render3d.Color {
	return render3d.NewColorRGB(176.0/255, 144.0/255, 26.0/255)
}

func nuggetOutline() *model2d.Mesh {
	bezier := model2d.BezierCurve{
		model2d.XY(0, 0),
		model2d.XY(0.2, 0.2),
		model2d.XY(0.5, 1.0),
		model2d.XY(0, 0.5),
	}
	res := model2d.NewMesh()
	for i := 0; i < 1000; i++ {
		t1 := float64(i) / 1000.0
		t2 := float64(i+1) / 1000.0
		res.Add(&model2d.Segment{
			bezier.Eval(t1),
			bezier.Eval(t2),
		})
	}
	res.AddMesh(res.MapCoords(model2d.XY(-1, 1).Mul))
	res, _ = res.Repair(1e-5).RepairNormals(1e-5)
	return res
}
