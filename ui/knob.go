package ui

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	KNOB_SIZE = image.Rect(0, 0, 100, 120)

	// from https://clrs.cc/
	black   = color.RGBA{0x11, 0x11, 0x11, 0xff}
	silver  = color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	fuschia = color.RGBA{0xf0, 0x12, 0xbe, 0xff}
	yellow  = color.RGBA{0xff, 0xDC, 0x00, 0xff}
	navy    = color.RGBA{0x00, 0x1F, 0x3F, 0xaa}
)

// draws a filled rounded rect
func boxBackground(gc *draw2dimg.GraphicContext, rect image.Rectangle) {
	pad := 2.0
	radius := 10.0

	pi := math.Pi
	halfpi := pi / 2

	minX := float64(rect.Min.X)
	maxX := float64(rect.Max.X)
	minY := float64(rect.Min.Y)
	maxY := float64(rect.Max.Y)

	gc.ArcTo(minX+radius+pad, minY+radius+pad, radius, radius, pi, halfpi)
	gc.LineTo(maxX-radius-pad, minY+pad)
	gc.ArcTo(maxX-radius-pad, minY+radius+pad, radius, radius, -halfpi, halfpi)
	gc.LineTo(maxX-pad, maxY-radius-pad)
	gc.ArcTo(maxX-radius-pad, maxY-radius-pad, radius, radius, 0, halfpi)
	gc.LineTo(minX+radius+pad, maxY-pad)
	gc.ArcTo(minX+radius+pad, maxY-radius-pad, radius, radius, halfpi, halfpi)
	gc.Close()
	gc.FillStroke()

}

type knob struct {
	draw.Image
	value float32
	label string
	gc    *draw2dimg.GraphicContext
}

func Knob(label string) *knob {
	k := &knob{label: label}
	k.Image = image.NewRGBA(KNOB_SIZE)
	k.gc = draw2dimg.NewGraphicContext(k.Image)

	k.gc.SetFillColor(navy)
	k.gc.SetStrokeColor(silver)
	k.gc.SetLineWidth(1)

	draw2dkit.RoundedRectangle(k.gc,
		float64(k.Bounds().Min.X+1), float64(k.Bounds().Min.Y+1),
		float64(k.Bounds().Max.X-1), float64(k.Bounds().Max.Y-1),
		10, 10)
	k.gc.FillStroke()

	center := k.Bounds().Size().Div(2)

	k.gc.SetFontData(draw2d.FontData{
		Name:   "SFAlienEncountersSolid",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleItalic,
	})
	k.gc.SetFontSize(15)
	k.gc.SetLineWidth(0.5)
	k.gc.SetStrokeColor(black)
	k.gc.SetFillColor(yellow)
	l, _, r, _ := k.gc.GetStringBounds(k.label)
	k.gc.FillStringAt(k.label, float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))
	k.gc.StrokeStringAt(k.label, float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))

	k.gc.SetStrokeColor(fuschia)
	k.gc.SetLineWidth(6)
	k.gc.SetLineCap(draw2d.RoundCap)
	startAngle := 110.0 * (math.Pi / 180.0)
	angle := 280.0 * (math.Pi / 180.0)
	k.gc.ArcTo(float64(center.X), float64(center.X), 35, 35, startAngle, angle)
	k.gc.Stroke()

	return k
}
