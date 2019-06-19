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
)

type knob struct {
	*pane
	param patch.Param
}

func Knob(bounds image.Rectangle, param patch.Param) *knob {
	k := &knob{
		pane:  BlankPane(bounds, PURPLE),
		param: param,
	}

	return k
}

func (k *knob) NeedsUpdate(id patch.ParamId) image.Rectangle {
	if k.param.ID() == id {
		return k.Bounds()
	}
	return image.ZR
}

func (k *knob) paint(bounds image.Rectangle) {
	if bounds == image.ZR {
		return
	}
	if bounds != k.Bounds() {
		bounds = k.Bounds()
		fmt.Printf("Warning: was asked to partial update a knob\n")
	}
	draw.Draw(k, k.Bounds(), image.NewUniform(color.Transparent), image.ZP, draw.Src)
	gc := draw2dimg.NewGraphicContext(k.Image)

	center := k.Bounds().Size().Div(2)

	gc.SetFontData(draw2d.FontData{
		Name:   "SFAlienEncountersSolid",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleItalic,
	})
	gc.SetFontSize(15)
	gc.SetLineWidth(0.5)
	gc.SetStrokeColor(BLACK)
	gc.SetFillColor(YELLOW)
	l, _, r, _ := gc.GetStringBounds(k.param.Label())
	gc.FillStringAt(k.param.Label(), float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))
	gc.StrokeStringAt(k.param.Label(), float64(center.X)-(r-l)/2.0, float64(k.Bounds().Max.Y-10))

	gc.SetStrokeColor(FUSCHIA)
	gc.SetLineWidth(6)
	gc.SetLineCap(draw2d.RoundCap)

	startAngle := 112.5 * (math.Pi / 180.0)
	angleTravel := 315.0 * (math.Pi / 180.0)
	angle := float64(k.param.ValAsCC()) / 128.0 * angleTravel
	//fmt.Printf("knob %s, val is %d, proportion %.2f, angle %.2f\n", k.param.Label(), k.param.ValAsCC(), float64(k.param.ValAsCC())/128.0, angle)

	gc.ArcTo(float64(center.X), float64(center.X), 35, 35, startAngle, angle)
	gc.Stroke()
}
