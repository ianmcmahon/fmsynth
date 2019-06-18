package ui

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/ianmcmahon/fmsynth/audio"
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

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 480
)

type screen struct {
	draw.Image
	bounds image.Rectangle
	layout *layout

	window wde.Window // osx specific; this needs to get factored out somehow
}

func (s *screen) paint(bounds image.Rectangle) {
	fmt.Printf("in screen paint: %v\n", bounds)
	s.layout.paint(bounds)
	draw.Draw(s.Image, bounds, s.layout, bounds.Min, draw.Src)
}

// temporary
var engine *audio.Engine

func Start(eng *audio.Engine) {
	engine = eng
	//// osx specific
	window, err := wde.NewWindow(SCREEN_WIDTH, SCREEN_HEIGHT)
	if err != nil {
		panic(err)
	}
	window.SetTitle("such synth")
	window.LockSize(true)
	/////

	draw2d.SetFontCache(draw2d.NewFolderFontCache("ui/fonts"))

	screenBounds := image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)

	screen := &screen{
		Image:  image.NewRGBA(screenBounds),
		bounds: screenBounds,
		window: window,
		layout: SampleLayout(screenBounds),
	}

	screen.paint(screenBounds)

	window.Screen().CopyRGBA(screen.Image.(*image.RGBA), screen.Bounds())
	window.FlushImage()
	window.Show()
}
