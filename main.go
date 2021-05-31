package main

import (
	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

func main() {
	penguin := NewPenguin()

	mesh := model3d.MarchingCubesSearch(penguin, 0.01, 8)
	SaveMesh(mesh, penguin)
	RenderMesh(mesh, penguin)
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
