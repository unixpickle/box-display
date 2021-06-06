package main

import (
	"math"

	"github.com/unixpickle/model3d/model2d"
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const (
	FlowerPedalThickness = 0.05
	FlowerPedalHeight    = 0.8
	FlowerStemRadius     = 0.1
	FlowerColorEpsilon   = 0.01
)

type Flower struct {
	model3d.Solid
	stemCheck model3d.Solid
}

func NewFlower() *Flower {
	var pedals model3d.JoinedSolid
	frequencies := []float64{5.0, 5.0, 5.0}
	phases := []float64{0.0, math.Pi * 2 / 10.0, math.Pi * 2 / 5.0}
	for i, length := range []float64{2.0, 4.0, 9.0} {
		fp := &flowerPedals{
			MinArcLength: length - 0.4,
			MaxArcLength: length + 0.4,
			Frequency:    frequencies[i],
			PhaseShift:   phases[i],
		}
		pedals = append(pedals, fp.Solid(0.5))
	}
	stem := model3d.Cylinder{
		P1:     model3d.Z(0.89),
		P2:     model3d.Z(1.55),
		Radius: FlowerStemRadius,
	}
	thorn1 := model3d.Cone{
		Base:   model3d.XZ(-0.05, 1.1),
		Tip:    model3d.XZ(-0.2, 1.1),
		Radius: 0.075,
	}
	thorn2 := model3d.Cone{
		Base:   model3d.XZ(0.05, 1.3),
		Tip:    model3d.XZ(0.2, 1.3),
		Radius: 0.075,
	}
	stemCheck := stem
	thorn1Check := thorn1
	thorn2Check := thorn2
	stemCheck.Radius += FlowerColorEpsilon
	thorn1Check.Radius += FlowerColorEpsilon
	thorn2Check.Radius += FlowerColorEpsilon
	return &Flower{
		Solid: model3d.JoinedSolid{
			model3d.TransformSolid(
				model3d.JoinedTransform{
					model3d.Rotation(model3d.XY(1, 0.2).Normalize(), 0.8),
					&model3d.Translate{Offset: model3d.Z(1.5)},
				},
				pedals,
			),
			&stem,
			&thorn1,
			&thorn2,
		},
		stemCheck: model3d.JoinedSolid{&stemCheck, &thorn1Check, &thorn2Check},
	}
}

func (f *Flower) Color(c model3d.Coord3D) render3d.Color {
	if f.stemCheck.Contains(c) {
		return render3d.NewColorRGB(0, 0.7, 0)
	}
	return render3d.NewColorRGB(1.0, 0, 0)
}

type flowerPedals struct {
	MinArcLength float64
	MaxArcLength float64
	Frequency    float64
	PhaseShift   float64
}

func (f *flowerPedals) arcLength(theta float64) float64 {
	frac := (math.Sin(f.Frequency*(theta+f.PhaseShift)) + 1) / 2
	return frac*(f.MaxArcLength-f.MinArcLength) + f.MinArcLength
}

func (f *flowerPedals) Mesh() *model3d.Mesh {
	// Create a "wiggly paraboloid" surface mesh by
	// iterating over theta, then iterating over the
	// arc of a 2D rotated parabola.
	res := model3d.NewMesh()
	divisions := 100
	if Production {
		divisions = 500
	}
	thetaForDiv := func(i int) float64 {
		if i < 0 {
			i += divisions
		}
		return math.Pi * 2 / float64(divisions) * float64(i%divisions)
	}
	radialPoint := func(theta float64, c2 model2d.Coord) model3d.Coord3D {
		return model3d.XYZ(math.Cos(theta)*c2.X, math.Sin(theta)*c2.X, c2.Y)
	}
	points, lengths := parabolicPoints(f.MaxArcLength)
	for i := 0; i < divisions; i++ {
		theta := thetaForDiv(i)
		nextTheta := thetaForDiv(i + 1)
		length := math.Min(f.arcLength(theta), f.arcLength(nextTheta))
		for j, l := range lengths {
			if l > length {
				break
			}
			if lengths[j+1] > length {
				// Create a quad which is tangent to the rim of
				// the wiggly paraboloid, by approximating the
				// final segment of each section using the slope
				// at points[j]->points[j+1].
				lastSeg := points[j+1].Sub(points[j]).Normalize()
				diff1 := f.arcLength(theta) - l
				diff2 := f.arcLength(nextTheta) - l
				end1 := points[j].Add(lastSeg.Scale(diff1))
				end2 := points[j].Add(lastSeg.Scale(diff2))
				res.AddQuad(
					radialPoint(theta, points[j]),
					radialPoint(nextTheta, points[j]),
					radialPoint(nextTheta, end2),
					radialPoint(theta, end1),
				)
			} else {
				res.AddQuad(
					radialPoint(theta, points[j]),
					radialPoint(nextTheta, points[j]),
					radialPoint(nextTheta, points[j+1]),
					radialPoint(theta, points[j+1]),
				)
			}
		}
	}
	return res
}

func (f *flowerPedals) Solid(desiredHeight float64) model3d.Solid {
	mesh := f.Mesh()
	mesh = mesh.Translate(model3d.Z(-mesh.Min().Z))
	height := mesh.Max().Z - mesh.Min().Z
	mesh = mesh.Scale(desiredHeight / height)
	return model3d.NewColliderSolidHollow(model3d.MeshToCollider(mesh), FlowerPedalThickness)
}

func parabolicPoints(maxArcLength float64) (points []model2d.Coord, lengths []float64) {
	dx := 0.01
	if Production {
		dx = 0.002
	}
	x := 0.0
	length := 0.0
	points = append(points, model2d.XY(0, 0))
	lengths = append(lengths, 0.0)
	for {
		prevY := x * x
		x += dx
		y := x * x
		length += math.Sqrt((prevY-y)*(prevY-y) + dx*dx)
		points = append(points, model2d.XY(x, y))
		lengths = append(lengths, length)
		if length > maxArcLength {
			break
		}
	}
	return
}
