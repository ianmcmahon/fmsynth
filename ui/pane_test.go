package ui

import (
	"image"
	"image/color"
	"testing"
)

type notifierPane struct {
	*pane
	lastUpdated image.Rectangle
}

func (p *notifierPane) paint(rect image.Rectangle) {
	p.lastUpdated = rect
	p.pane.paint(rect)
}

func TestChildRectangles(t *testing.T) {
	container := BlankPane(image.Rect(0, 0, 40, 40), color.Black)

	p1 := &notifierPane{BlankPane(image.Rect(0, 0, 10, 10), color.White), image.ZR}
	p2 := &notifierPane{BlankPane(image.Rect(0, 0, 20, 20), color.White), image.ZR}

	container.AddChild(p1, image.Pt(5, 5))
	container.AddChild(p2, image.Pt(15, 15))

	// paint the whole container and verify the images got the updates
	container.paint(container.Bounds())

	expectRect(t, p1.lastUpdated, p1.Bounds())
	expectRect(t, p2.lastUpdated, p2.Bounds())

	expectColor(t, container.At(2, 2), color.RGBAModel.Convert(color.Black))
	expectColor(t, container.At(12, 12), color.RGBAModel.Convert(color.White))
	expectColor(t, container.At(17, 17), color.RGBAModel.Convert(color.White))
	expectColor(t, container.At(34, 34), color.RGBAModel.Convert(color.White))
	expectColor(t, container.At(36, 36), color.RGBAModel.Convert(color.Black))

	/*
		before, err := os.Create("before.png")
		if err != nil {
			t.Fatal(err)
		}

		if err := png.Encode(before, container); err != nil {
			t.Fatal(err)
		}
	*/

	p1.pane.backgroundColor = color.RGBA{0xff, 0x00, 0x00, 0xFF}
	p2.pane.backgroundColor = color.RGBA{0x00, 0xff, 0x00, 0xFF}
	container.backgroundColor = color.RGBA{0x00, 0x00, 0xff, 0xFF}
	container.paint(image.Rect(10, 10, 20, 20))

	/*
		after, err := os.Create("after.png")
		if err != nil {
			t.Fatal(err)
		}

		if err := png.Encode(after, container); err != nil {
			t.Fatal(err)
		}
	*/

	expectColor(t, container.At(2, 2), color.RGBAModel.Convert(color.Black))
	expectColor(t, container.At(12, 12), color.RGBA{0xff, 0x00, 0x00, 0xff})
	expectColor(t, container.At(17, 17), color.RGBA{0x00, 0xff, 0x00, 0xff})
	expectColor(t, container.At(14, 17), color.RGBA{0x00, 0x00, 0xff, 0xff})
	expectColor(t, container.At(34, 34), color.RGBAModel.Convert(color.White))
	expectColor(t, container.At(36, 36), color.RGBAModel.Convert(color.Black))

	expectRect(t, p1.lastUpdated, image.Rect(5, 5, 10, 10))
	expectRect(t, p2.lastUpdated, image.Rect(0, 0, 5, 5))
}

func expectRect(t *testing.T, a, b image.Rectangle) {
	if a != b {
		t.Errorf("expected %#v, got %#v\n", b, a)
	}
}

func expectColor(t *testing.T, a, b color.Color) {
	if a != b {
		t.Errorf("expected %#v, got %#v\n", b, a)
	}
}
