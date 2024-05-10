package main

import (
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"strconv"
)

type Pixel struct {
	*canvas.Rectangle
	round int
}

func (b *Board) String() string {
	return strconv.Itoa(b.round)
}

func (p *Pixel) IsWhite() bool {
	return p.FillColor == color.White
}

func (p *Pixel) SetColor(c color.Color) {
	p.FillColor = c
	p.Refresh()
}
