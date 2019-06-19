package ui

import (
	"image"

	"github.com/ianmcmahon/fmsynth/patch"
)

type paramPage struct {
	Pane
	params []patch.Param
}

func ParamPage(bounds image.Rectangle, params []patch.Param) *paramPage {
	page := &paramPage{
		RoundedRectPane(bounds),
		params,
	}

	for i, p := range params {
		childBounds := page.controlBounds(i)
		page.AddChild(Knob(childBounds, p), childBounds.Min)
	}

	return page
}

func (p *paramPage) controlSize() image.Rectangle {
	return image.Rect(0, 0, p.Bounds().Dx()/4, p.Bounds().Dy()/2)
}

func (p *paramPage) controlBounds(i int) image.Rectangle {
	size := p.controlSize()
	at := image.Pt(i%4*size.Dx(), i/4*size.Dy())
	return size.Add(at)
}

/*
func (p *paramPage) NeedsUpdate(id patch.ParamId) image.Rectangle {
	rect := image.ZR
	for _, c := range p.children {
		rect = rect.Union(c.NeedsUpdate(id).Add(c.at))
	}
	fmt.Printf("param page needs update: %v\n", rect)
	return rect
}
*/
