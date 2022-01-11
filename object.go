package main

import (
	"math"

	"github.com/unixpickle/model3d/model3d"
	"github.com/unixpickle/model3d/render3d"
	"github.com/unixpickle/model3d/toolbox3d"
)

const (
	joinColorPrecDev  = 0.02
	joinColorPrecProd = 0.005
)

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

func Scale(obj Object, s float64) Object {
	return Transform(&model3d.Scale{Scale: s}, obj)
}

func Translate(obj Object, offset model3d.Coord3D) Object {
	return Transform(&model3d.Translate{Offset: offset}, obj)
}

func RotateZ(obj Object, angle float64) Object {
	return Transform(model3d.Rotation(model3d.Z(1), angle), obj)
}

func (t *transformedObject) Color(c model3d.Coord3D) render3d.Color {
	return t.obj.Color(t.invTransform.Apply(c))
}

type joinedObject struct {
	model3d.Solid
	bounds  []*model3d.Rect
	colorFn toolbox3d.CoordColorFunc
}

func Join(objs ...Object) Object {
	js := make(model3d.JoinedSolid, len(objs))
	bounds := make([]*model3d.Rect, len(objs))
	var colorArgs []interface{}
	for i, obj := range objs {
		js[i] = obj
		size := obj.Max().Sub(obj.Min())
		prec := joinColorPrecDev
		if Production {
			prec = joinColorPrecProd
		}
		eps := math.Max(size.X, math.Max(size.Y, size.Z)) * prec
		mesh := model3d.MarchingCubesSearch(obj, eps, 8)
		sdf := model3d.MeshToSDF(mesh)
		colorArgs = append(colorArgs, sdf, obj.Color)
		bounds[i] = model3d.NewRect(mesh.Min(), mesh.Max())
	}
	return &joinedObject{
		Solid:   js,
		bounds:  bounds,
		colorFn: toolbox3d.JoinedCoordColorFunc(colorArgs...),
	}
}

func (j *joinedObject) Color(c model3d.Coord3D) render3d.Color {
	return j.colorFn(c)
}
