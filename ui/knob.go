package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/ianmcmahon/fmsynth/patch"
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

type control interface {
	draw.Image
	paint(image.Rectangle)
}

type knob struct {
	draw.Image
	param patch.Param
	gc    *draw2dimg.GraphicContext
}

func Knob(bounds image.Rectangle, param patch.Param) *knob {
	k := &knob{
		Image: image.NewRGBA(bounds),
		param: param,
	}
	k.gc = draw2dimg.NewGraphicContext(k.Image)

	return k
}

func (k *knob) paint(bounds image.Rectangle) {
	fmt.Printf("in paint knob %s  %v\n", k.param.Label(), bounds)
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
	l, _, r, _ := k.gc.GetStringBounds(k.param.Label())
	k.gc.FillStringAt(k.param.Label(), float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))
	k.gc.StrokeStringAt(k.param.Label(), float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))

	k.gc.SetStrokeColor(fuschia)
	k.gc.SetLineWidth(6)
	k.gc.SetLineCap(draw2d.RoundCap)
	startAngle := 110.0 * (math.Pi / 180.0)
	angle := 280.0 * (math.Pi / 180.0)
	k.gc.ArcTo(float64(center.X), float64(center.X), 35, 35, startAngle, angle)
	k.gc.Stroke()
}
