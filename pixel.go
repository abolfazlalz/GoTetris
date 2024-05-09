package main

import (
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

type Pixel struct {
	*canvas.Rectangle
	round int
}

func (p *Pixel) IsWhite() bool {
	return p.FillColor == color.White
}

func (p *Pixel) SetColor(c color.Color) {
	p.FillColor = c
	p.Refresh()
}
