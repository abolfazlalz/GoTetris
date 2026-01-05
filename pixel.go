package main

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2/canvas"
)

type Pixel struct {
	*canvas.Rectangle
	round int
}

func (b *Board) String() string {
	return strconv.Itoa(0)
}

func (p *Pixel) IsWhite() bool {
	return p.FillColor == color.White
}

func (p *Pixel) SetColor(c color.Color) {
	p.FillColor = c
	p.Refresh()
}
