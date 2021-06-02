package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
)

const joinColorPrec = 0.01

type Object interface {
	model3d.Solid

	Color(c model3d.Coord3D) render3d.Color
}

type transformedObject struct {
	model3d.Solid
	obj          Object
	invTransform model3d.Transform
}

func Transform(t model3d.Transform, obj Object) Object {
	return &transformedObject{
		Solid:        model3d.TransformSolid(t, obj),
		obj:          obj,
		invTransform: t.Inverse(),
	}
}

func Translate(obj Object, offset model3d.Coord3D) Object {
	return Transform(&model3d.Translate{Offset: offset}, obj)
}

func (t *transformedObject) Color(c model3d.Coord3D) render3d.Color {
	return t.obj.Color(t.invTransform.Apply(c))
}

type joinedObject struct {
	model3d.Solid
	objects []Object
	sdfs    []model3d.SDF
}

func Join(objs ...Object) Object {
	js := make(model3d.JoinedSolid, len(objs))
	sdfs := make([]model3d.SDF, len(objs))
	for i, obj := range objs {
		js[i] = obj
		size := obj.Max().Sub(obj.Min())
		eps := math.Max(size.X, math.Max(size.Y, size.Z)) * joinColorPrec
		mesh := model3d.MarchingCubesSearch(obj, eps, 8)
		sdfs[i] = model3d.MeshToSDF(mesh)
	}
	return &joinedObject{
		Solid:   js,
		objects: objs,
		sdfs:    sdfs,
	}
}

func (j *joinedObject) Color(c model3d.Coord3D) render3d.Color {
	maxSDF := math.Inf(-1)
	var closest Object
	for i, obj := range j.objects {
		sdf := j.sdfs[i].SDF(c)
		if sdf > maxSDF {
			maxSDF = sdf
			closest = obj
		}
	}
	return closest.Color(c)
}
