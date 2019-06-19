package ui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ianmcmahon/fmsynth/patch"
)

type Pane interface {
	image.Image
	paint(image.Rectangle)
	AddChild(c Pane, at image.Point)
	NeedsUpdate(id patch.ParamId) image.Rectangle
}

type child struct {
	Pane
	bounds image.Rectangle
}

type pane struct {
	draw.Image
	backgroundColor color.Color
	children        []child
}

func BlankPane(size image.Rectangle, bgColor color.Color) *pane {
	return &pane{
		Image:           image.NewRGBA(size.Sub(size.Min)),
		backgroundColor: bgColor,
		children:        make([]child, 0),
	}
}

func (p *pane) AddChild(c Pane, at image.Point) {
	p.children = append(p.children, child{c, c.Bounds().Add(at)})
}

func (p *pane) NeedsUpdate(id patch.ParamId) image.Rectangle {
	rect := image.ZR
	for _, child := range p.children {
		rect = rect.Union(child.NeedsUpdate(id).Add(child.bounds.Min))
	}
	return rect
}

func (parent *pane) UpdateChildren(rect image.Rectangle) {
	for _, child := range parent.children {
		// parentRect is the area of the parent image that overlaps this child
		parentRect := rect.Intersect(child.bounds)
		if parentRect != image.ZR {
			// childRect is the area in the child's coordinate space that needs painting
			childRect := parentRect.Sub(child.bounds.Min)
			child.paint(childRect)
			draw.Draw(parent, parentRect, child, childRect.Min, draw.Over)
		}
	}
}

func (parent *pane) paint(rect image.Rectangle) {
	draw.Draw(parent, rect, image.NewUniform(parent.backgroundColor), image.ZP, draw.Src)
	parent.UpdateChildren(rect)
}
