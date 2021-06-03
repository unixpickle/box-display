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
		Translate(NewPenguin(), model3d.X(-1)),
		NewFrame(),
		Translate(NewNugget(), model3d.X(1)),
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
	render3d.SaveRandomGrid("rendering.png", mesh, 3, 3, 400, colorFunc)
}
