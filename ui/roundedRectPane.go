package ui

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type roundedRectPane struct {
	*pane
	radius    int
	padding   int
	lineWidth int
	fill      color.Color
	stroke    color.Color
}

func RoundedRectPane(bounds image.Rectangle) *roundedRectPane {
	return &roundedRectPane{
		pane:      BlankPane(bounds, color.Transparent),
		radius:    5,
		padding:   2,
		lineWidth: 1,
		fill:      NAVY,
		stroke:    SILVER,
	}
}

func (p *roundedRectPane) paint(rect image.Rectangle) {
	// clear the canvas to transparent
	//	draw.Draw(p, rect, image.NewUniform(color.Transparent), image.ZP, draw.Src)

	// on updates this will stroke the whole rectangle, this is fine
	// we clear the portion of image that we're going to redraw, and re-filling the whole
	// rect ensures that the background under children is redrawn properly
	gc := draw2dimg.NewGraphicContext(p.Image)
	gc.SetFillColor(p.fill)
	gc.SetStrokeColor(p.stroke)
	gc.SetLineWidth(float64(p.lineWidth))

	b := p.Bounds().Inset(p.padding)
	draw2dkit.RoundedRectangle(gc,
		float64(b.Min.X), float64(b.Min.Y),
		float64(b.Max.X), float64(b.Max.Y),
		float64(p.radius), float64(p.radius))
	gc.FillStroke()

	p.UpdateChildren(rect)
}
