package ui

import (
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
	Pane
	visible bool
}

func SampleLayout(bounds image.Rectangle) *layout {
	backgroundFile, err := os.Open("ui/backgrounds/efe-kurnaz-315384-unsplash.jpg")
	if err != nil {
		panic(err)
	}

	backgroundImg, err := jpeg.Decode(backgroundFile)
	if err != nil {
		panic(err)
	}

	layout := &layout{Pane: BackgroundPane(bounds, backgroundImg)}

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
	paramPage := ParamPage(image.Rect(0, 0, 400, 240), params)
	layout.AddChild(paramPage, image.Pt(20, 20))

	return layout
}

/*
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
*/
