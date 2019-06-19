package ui

import (
	"image"
	"image/draw"
)

type backgroundPane struct {
	*pane
	backgroundImage image.Image
}

func BackgroundPane(bounds image.Rectangle, img image.Image) *backgroundPane {
	return &backgroundPane{BlankPane(bounds, BLACK), img}
}

func (p *backgroundPane) paint(rect image.Rectangle) {
	draw.Draw(p, rect, p.backgroundImage, image.ZP, draw.Src)
	p.UpdateChildren(rect)
}
