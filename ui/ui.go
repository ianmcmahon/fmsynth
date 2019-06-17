package ui

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/llgcode/draw2d"
	wde "github.com/skelterjohn/go.wde"
)

/*
	UI Concepts

	'screen' is the base canvas of the display
	'layout' manages the background image and placement of child containers
	'param page' contains eight individual param controls that map to hardware encoders
	  as well as any background drawing it may do
	a param control such as 'knob' knows how to draw its labels, indicators, and controls
	  as well as any background drawing it may do
	a Patch is a data structure that contains all of the params for the engine
	a Param contains the concrete value of record, as well as metadata about how it
	  maps to synth engine components, midi CCs, and UI elements, including label and font
		for theming
	param pages will be configured by paramId, so a given page can ask the patch for param
	  metadata
	it'd be possible to have user configurable param pages, or maybe just a macro page and
	  a way to edit macros binding arbitrary parameters to them with attenuverting that all sum
*/

func Start() {
	window, err := wde.NewWindow(800, 480)
	if err != nil {
		panic(err)
	}

	window.SetTitle("such synth")
	window.LockSize(true)
	window.Show()

	draw2d.SetFontCache(draw2d.NewFolderFontCache("ui/fonts"))

	screen := image.NewRGBA(image.Rect(0, 0, 800, 480))

	backgroundFile, err := os.Open("ui/backgrounds/efe-kurnaz-315384-unsplash.jpg")
	if err != nil {
		panic(err)
	}

	backgroundImg, err := jpeg.Decode(backgroundFile)
	if err != nil {
		panic(err)
	}

	draw.Draw(screen, screen.Bounds(), backgroundImg, image.ZP, draw.Src)

	page := image.NewRGBA(image.Rect(0, 0, KNOB_SIZE.Max.X*4, KNOB_SIZE.Max.Y*2))

	knobs := []string{"ATTACC", "DECAY", "SUSTN", "RELEASE", "ALG", "FEEDBK", "MIX", "BKGRD"}
	xSize := KNOB_SIZE.Max.X
	ySize := KNOB_SIZE.Max.Y
	for i := 0; i < 8; i++ {
		k := Knob(knobs[i])
		draw.Draw(page, k.Bounds().Add(image.Pt(i%4*xSize, i/4*ySize)), k, image.ZP, draw.Over)
	}

	draw.Draw(screen, page.Bounds().Add(image.Pt(20, 20)), page, image.ZP, draw.Over)

	window.Screen().CopyRGBA(screen, screen.Bounds())
	window.FlushImage()

}
