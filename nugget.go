package main

import (
	"math"
	"math/rand"

	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	NuggetRounding     = 0.1
	NuggetCrumbRadius  = 0.03
	NuggetCrumbEpsilon = 0.01
	NuggetCrumbMinDist = NuggetCrumbRadius / 2
)

type Nugget struct {
	mainBody  model3d.SDF
	centers   *model3d.CoordTree
	colors    map[model3d.Coord3D]render3d.Color
	baseColor render3d.Color
}

func NewNugget() *Nugget {
	outline := model2d.MeshToSDF(nuggetOutline())

	mainBody := model3d.ProfileSDF(outline, -0.1, 0.1)
	mainBody = model3d.TransformSDF(model3d.Rotation(model3d.X(1), math.Pi/2), mainBody)

	baseColor := render3d.NewColorRGB(176.0/255, 144.0/255, 26.0/255)
	allColors := nuggetColors(baseColor)

	// Sample points approximately on the surface of the solid.
	roughMesh := model3d.MarchingCubesSearch(&Nugget{mainBody: mainBody}, 0.03, 8)
	light := render3d.NewMeshAreaLight(roughMesh, render3d.Color{})
	gen := rand.New(rand.NewSource(0))
	colorMapping := map[model3d.Coord3D]render3d.Color{}
	var centers []model3d.Coord3D
	for i := 0; i < 3000; i++ {
		var point model3d.Coord3D
	SampleLoop:
		for {
			point, _, _ = light.SampleLight(gen)
			for _, p2 := range centers {
				if point.Dist(p2) < NuggetCrumbMinDist {
					continue SampleLoop
				}
			}
			break
		}
		centers = append(centers, point)
		colorMapping[point] = allColors[rand.Intn(len(allColors))]
	}

	return &Nugget{
		mainBody:  mainBody,
		centers:   model3d.NewCoordTree(centers),
		colors:    colorMapping,
		baseColor: baseColor,
	}
}

func (n *Nugget) Min() model3d.Coord3D {
	return n.mainBody.Min().Sub(model3d.XYZ(1, 1, 1).Scale(NuggetCrumbRadius + NuggetRounding))
}

func (n *Nugget) Max() model3d.Coord3D {
	return n.mainBody.Max().Add(model3d.XYZ(1, 1, 1).Scale(NuggetCrumbRadius + NuggetRounding))
}

func (n *Nugget) Contains(c model3d.Coord3D) bool {
	if n.mainBody.SDF(c) > -NuggetRounding {
		return true
	}
	if n.centers != nil {
		return n.centers.SphereCollision(c, NuggetCrumbRadius)
	}
	return false
}

func (n *Nugget) Color(c model3d.Coord3D) render3d.Color {
	if n.centers != nil {
		neighbor := n.centers.NearestNeighbor(c)
		if c.Dist(neighbor) < NuggetCrumbRadius+NuggetCrumbEpsilon {
			return n.colors[neighbor]
		}
	}
	return n.baseColor
}

func nuggetColors(base render3d.Color) []render3d.Color {
	var res []render3d.Color
	for i := -0.1; i <= 0.1; i += 0.025 {
		res = append(res, base.Scale(1+i))
	}
	return res
}

func nuggetOutline() *model2d.Mesh {
	bezier := model2d.BezierCurve{
		model2d.XY(0, 0),
		model2d.XY(0.3, 0.2),
		model2d.XY(0.6, 1.0),
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
