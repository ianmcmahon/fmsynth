package ui

import (
	"image"
	"image/draw"

	wde "github.com/skelterjohn/go.wde"
)

func Start() {
	window, err := wde.NewWindow(800, 480)
	if err != nil {
		panic(err)
	}

	window.SetTitle("such synth")
	window.LockSize(true)
	window.Show()

	page := image.NewRGBA(image.Rect(0, 0, 800, 480))
	addLabel(page, 500, 300, "do words to it")

	k := Knob("Turn Me")
	draw.Draw(page, k.Bounds().Add(image.Pt(10, 10)), k, image.ZP, draw.Src)

	window.Screen().CopyRGBA(page, page.Bounds())
	window.FlushImage()

}
