package ui

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/ianmcmahon/fmsynth/patch"
)

type layout struct {
	draw.Image
	visible    bool
	background image.Image
	paramPage  *paramPage
}

func SampleLayout(bounds image.Rectangle) *layout {
	layout := &layout{
		Image: image.NewRGBA(bounds),
	}

	backgroundFile, err := os.Open("ui/backgrounds/efe-kurnaz-315384-unsplash.jpg")
	if err != nil {
		panic(err)
	}

	layout.background, err = jpeg.Decode(backgroundFile)
	if err != nil {
		panic(err)
	}

	ptch := engine.CurrentPatch()
	params := []patch.Param{
		ptch.GetParam(patch.PATCH_ALGORITHM),
		ptch.GetParam(patch.OPR_RATIO | patch.GRP_A),
		ptch.GetParam(patch.OPR_RATIO | patch.GRP_B1),
		ptch.GetParam(patch.OPR_RATIO | patch.GRP_C),
		ptch.GetParam(patch.ENV_ATTACK | patch.GRP_VCA),
		ptch.GetParam(patch.ENV_DECAY | patch.GRP_VCA),
		ptch.GetParam(patch.ENV_SUSTAIN | patch.GRP_VCA),
		ptch.GetParam(patch.ENV_RELEASE | patch.GRP_VCA),
	}

	layout.paramPage = newParamPage(params)

	return layout
}

func (l *layout) paint(bounds image.Rectangle) {
	fmt.Printf("in layout paint %v\n", bounds)
	draw.Draw(l.Image, bounds, l.background, bounds.Min, draw.Src)

	// todo: check if bounds intersects the page
	l.paramPage.paint(l.paramPage.Bounds()) // hacky, just forces repaint of the whole container

	draw.Draw(l.Image, l.paramPage.Bounds(), l.paramPage, image.ZP, draw.Over)
}

type paramPage struct {
	draw.Image
	params   []patch.Param
	controls []control
	visible  bool
	active   bool
}

// this very much assumes 8 params, 2 rows of 4, 400x240 (gives 100x120 tiles)
func newParamPage(params []patch.Param) *paramPage {
	page := &paramPage{
		Image:    image.NewRGBA(image.Rect(0, 0, 400, 240)),
		params:   params,
		controls: make([]control, len(params)),
	}

	for i, p := range params {
		// todo: introspect the param to decide what control to use
		b := image.Rect(0, 0, page.Bounds().Dx()/4, page.Bounds().Dy()/2)
		page.controls[i] = Knob(b, p)
	}

	return page
}

func (p *paramPage) paint(bounds image.Rectangle) {
	fmt.Printf("in paramPage paint: %v\n", bounds)
	for i, control := range p.controls {
		// todo: check if bounds intersects the param
		// and not everything will be a knob!
		control.paint(control.Bounds())
		xSize := control.Bounds().Dx()
		ySize := control.Bounds().Dy()
		destRect := control.Bounds().Add(image.Pt(i%4*xSize, i/4*ySize))
		fmt.Printf("control %d painting at %v\n", i, destRect)
		draw.Draw(p, destRect, control, image.ZP, draw.Over)
	}
}
