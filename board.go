package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
	"sync"
	"time"
)

type Board struct {
	width    int
	height   int
	perPixel int
	rowPixel int
	colPixel int
	round    int
	mu       sync.Mutex
	pixels   [][]*Pixel
}

func NewBoard(width, height, perPixel int) *Board {

	return &Board{
		width:    width,
		height:   height,
		perPixel: perPixel,
		rowPixel: width / perPixel,
		colPixel: height / perPixel,
		pixels:   make([][]*Pixel, height/perPixel),
		round:    1,
		mu:       sync.Mutex{},
	}
}

func (b *Board) checkCollisionPixel(col, row int) (right, left, bottom bool) {
	left = row == 0
	bottom = col >= b.colPixel-1
	right = row == b.rowPixel-1
	return
}

func (b *Board) checkCollisionRound() (right, left, bottom bool) {
	right = true
	left = true
	bottom = true
	for i := 0; i < b.colPixel; i++ {
		for j := 0; j < b.rowPixel; j++ {
			pixel := b.pixels[i][j]
			if pixel.IsWhite() || pixel.round != b.round {
				continue
			}
			newRight, newLeft, newBottom := b.checkCollisionPixel(i, j)
			right = !newRight && right
			left = !newLeft && left
			bottom = !newBottom && bottom
		}
	}
	return !right, !left, !bottom
}

func (b *Board) Left() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, left, _ := b.checkCollisionRound(); left {
		return
	}
	for i := b.colPixel - 1; 0 < i; i-- {
		for j := 1; j < b.rowPixel; j++ {
			pixel := b.pixels[i][j]
			if pixel.IsWhite() || pixel.round != b.round {
				continue
			}
			b.pixels[i][j-1].SetColor(pixel.FillColor)
			pixel.SetColor(color.White)
		}
	}
}

func (b *Board) Right() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if right, _, _ := b.checkCollisionRound(); right {
		return
	}
	for i := 0; i < b.colPixel; i++ {
		for j := b.rowPixel - 1; 0 <= j; j-- {
			pixel := b.pixels[i][j]
			if pixel.IsWhite() || pixel.round != b.round {
				continue
			}
			b.pixels[i][j+1].SetColor(pixel.FillColor)
			pixel.SetColor(color.White)
		}
	}
}

func (b *Board) Down() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, _, bottom := b.checkCollisionRound(); bottom {
		return false
	}
	for i := b.colPixel - 1; 0 <= i; i-- {
		for j := 0; j < b.rowPixel; j++ {
			pixel := b.pixels[i][j]
			if pixel.IsWhite() || pixel.round != b.round {
				continue
			}
			b.pixels[i+1][j].SetColor(pixel.FillColor)
			pixel.SetColor(color.White)
		}
	}
	return true
}

func (b *Board) Animate() {
	b.pixels[0][1].SetColor(color.RGBA{A: 255, R: 255})
	for range time.Tick(time.Millisecond * 550) {
		if x := b.Down(); !x {
		}
	}
}

func (b *Board) SetWindowContent(w fyne.Window) {
	items := make([]fyne.CanvasObject, 0)
	for i := 0; i < b.colPixel; i++ {
		b.pixels[i] = make([]*Pixel, b.rowPixel)
		for j := 0; j < b.rowPixel; j++ {
			b.pixels[i][j] = &Pixel{Rectangle: canvas.NewRectangle(color.White), round: 1}
			items = append(items, b.pixels[i][j].Rectangle)
		}
	}

	element := container.New(layout.NewGridLayout(b.rowPixel), items...)
	w.SetContent(element)
}
