package ui

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/ianmcmahon/fmsynth/patch"
)

type uicomponent interface {
	draw.Image
	NeedsUpdate(id patch.ParamId) image.Rectangle
	paint(image.Rectangle)
}

type placement struct {
	uicomponent
	at image.Point
}

type layout struct {
	draw.Image
	visible    bool
	background image.Image
	children   []placement
}

func SampleLayout(bounds image.Rectangle) *layout {
	layout := &layout{
		Image:    image.NewRGBA(bounds),
		children: make([]placement, 1),
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

	layout.children[0] = placement{newParamPage(params), image.Pt(20, 20)}

	return layout
}

func (l *layout) NeedsUpdate(id patch.ParamId) image.Rectangle {
	rect := image.ZR
	for _, child := range l.children {
		rect = rect.Union(child.NeedsUpdate(id))
	}
	return rect
}

func (l *layout) paint(bounds image.Rectangle) {
	if bounds == image.ZR {
		return
	}
	fmt.Printf("in layout paint %v\n", bounds)
	draw.Draw(l.Image, bounds, l.background, bounds.Min, draw.Src)

	for _, child := range l.children {
		// todo: check if bounds intersects the page
		childBounds := bounds.Intersect(child.Bounds().Add(child.at))
		fmt.Printf("layout.paint(%v) parent bounds: %v  child: %v bounds %v\n", bounds, childBounds.Add(child.at), child, childBounds)
		child.paint(childBounds)
		draw.Draw(l.Image, childBounds.Add(child.at), child, childBounds.Min, draw.Over)
	}

}

type paramPage struct {
	draw.Image
	params   []patch.Param
	children []placement
	visible  bool
	active   bool
}

// this very much assumes 8 params, 2 rows of 4, 400x240 (gives 100x120 tiles)
func newParamPage(params []patch.Param) *paramPage {
	page := &paramPage{
		Image:    image.NewRGBA(image.Rect(0, 0, 400, 240)),
		params:   params,
		children: make([]placement, len(params)),
	}

	for i, p := range params {
		// todo: introspect the param to decide what control to use
		page.children[i] = placement{Knob(page.controlSize(), p), page.controlBounds(i).Min}
		fmt.Printf("placing child %v at %v\n", page.children[i], page.children[i].at)
	}

	return page
}

func (p *paramPage) NeedsUpdate(id patch.ParamId) image.Rectangle {
	rect := image.ZR
	for _, c := range p.children {
		rect = rect.Union(c.NeedsUpdate(id).Add(c.at))
	}
	fmt.Printf("param page needs update: %v\n", rect)
	return rect
}

func (p *paramPage) controlSize() image.Rectangle {
	return image.Rect(0, 0, p.Bounds().Dx()/4, p.Bounds().Dy()/2)
}

func (p *paramPage) controlBounds(i int) image.Rectangle {
	size := p.controlSize()
	at := image.Pt(i%4*size.Dx(), i/4*size.Dy())
	return size.Add(at)
}

func (p *paramPage) paint(bounds image.Rectangle) {
	if bounds == image.ZR {
		return
	}
	fmt.Printf("in paramPage paint: %v\n", bounds)
	for i, child := range p.children {
		if bounds.Overlaps(child.Bounds().Add(child.at)) {
			child.paint(child.Bounds())
			draw.Draw(p, child.Bounds().Add(child.at), child, image.ZP, draw.Over)
			fmt.Printf("control %d painting at %v\n", i, child.Bounds().Add(child.at))
		}
	}
}
