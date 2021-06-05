package main

import (
	"log"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const Production = false

func main() {
	log.Println("Creating object...")
	obj := Join(
		Translate(RotateZ(NewPenguin(), 0.1), model3d.X(-1.3)),
		NewFrame(),
		Translate(
			RotateZ(Scale(Translate(NewNugget(), model3d.Z(0.05)), 1.4), -0.1),
			model3d.X(1.4),
		),
		Translate(NewSun(), model3d.XYZ(1.5, 0.6, 3.0)),
		NewVase(),
		Translate(NewCloud("models/cloud/cloud2.png"), model3d.XYZ(-1.5, 0.6, 3.2)),
		Translate(NewCloud("models/cloud/cloud1.png"), model3d.XYZ(0.1, 0.6, 2.7)),
	)

	log.Println("Creating mesh...")
	eps := 0.03
	if Production {
		eps = 0.01
	}
	mesh := model3d.MarchingCubesSearch(obj, eps, 8)

	log.Println("Rendering...")
	RenderMesh(mesh, obj)

	log.Println("Saving...")
	SaveMesh(mesh, obj)
}

func SaveMesh(mesh *model3d.Mesh, o Object) {
	vertexColor := func(c model3d.Coord3D) [3]float64 {
		r, g, b := render3d.RGB(o.Color(c))
		return [3]float64{r, g, b}
	}
	mesh.SaveMaterialOBJ("export.zip", model3d.VertexColorsToTriangle(vertexColor))
}

func RenderMesh(mesh *model3d.Mesh, o Object) {
	colorFunc := func(c model3d.Coord3D, rc model3d.RayCollision) render3d.Color {
		return o.Color(c)
	}
	res := 512
	if Production {
		res = 1024
	}
	render3d.SaveRendering(
		"rendering1.png",
		mesh,
		model3d.XYZ(2.5, -8.0, 2.0),
		res,
		res,
		colorFunc,
	)
	render3d.SaveRendering(
		"rendering2.png",
		mesh,
		model3d.XYZ(-2.5, -8.0, 3.0),
		res,
		res,
		colorFunc,
	)
}
