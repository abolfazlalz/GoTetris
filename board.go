package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

type Result int

const (
	GameOverResult Result = iota
	AddScoreResult
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
	colors   []color.Color
	score    int
	isDone   bool
	resultCh chan Result
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
		isDone:   false,
		resultCh: make(chan Result),
		score:    0,
		colors: []color.Color{
			color.RGBA{A: 255, R: 255},
			color.RGBA{A: 255, G: 255},
			color.RGBA{A: 255, B: 255},
			color.RGBA{A: 255, R: 255, B: 255},
			color.RGBA{A: 255, R: 255, G: 255},
			color.RGBA{A: 255, B: 255, G: 255},
		},
	}
}

func (b *Board) AddScore() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.score += 1
	b.resultCh <- AddScoreResult
}

func (b *Board) GameOver() {
	b.resultCh <- GameOverResult
	b.isDone = true
}

func (b *Board) checkCollisionPixel(col, row int) (right, left, bottom bool) {
	left = row == 0
	bottom = col >= b.colPixel-1
	right = row == b.rowPixel-1

	if !bottom {
		pixel := b.pixels[col+1][row]
		bottom = !pixel.IsWhite() && pixel.round != b.pixels[col][row].round
	}
	if !left {
		pixel := b.pixels[col][row-1]
		left = !pixel.IsWhite() && pixel.round != b.pixels[col][row].round
	}

	if !right {
		pixel := b.pixels[col][row+1]
		right = !pixel.IsWhite() && pixel.round != b.pixels[col][row].round
	}

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

func (b *Board) randomColor() color.Color {
	cLen := len(b.colors)
	return b.colors[rand.Intn(cLen)]
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
			b.pixels[i][j-1].round = pixel.round
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
			b.pixels[i][j+1].round = pixel.round
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
			b.pixels[i+1][j].round = pixel.round
			pixel.SetColor(color.White)
		}
	}
	return true
}

func (b *Board) checkForGameOver() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := 0; i < b.rowPixel; i++ {
		block := b.pixels[0][i]
		if !block.IsWhite() {
			b.GameOver()
			fmt.Println("game over !")
			continue
		}
	}
}

func (b *Board) checkForExplodeRow() {
	for i := 0; i < b.colPixel; i++ {
		isFilled := true
		for j := 0; j < b.rowPixel; j++ {
			isFilled = !b.pixels[i][j].IsWhite() && isFilled
		}
		if isFilled {
			b.AddScore()
			for j := 0; j < b.rowPixel; j++ {
				b.pixels[i][j].SetColor(color.White)
			}
			for i1 := i; 0 < i1; i1-- {
				for j1 := 0; j1 < b.rowPixel; j1++ {
					b.pixels[i1][j1].SetColor(b.pixels[i1-1][j1].FillColor)
				}
			}
		}
	}
}

func (b *Board) shapeGenerator() {
	b.mu.Lock()
	defer b.mu.Unlock()
	startXPoint := rand.Intn(b.rowPixel)
	nX := b.rowPixel - startXPoint
	if nX > 4 {
		nX = 4
	}
	endXPoint := startXPoint + rand.Intn(nX)

	startYPoint := 0
	nY := b.colPixel - startYPoint
	if nY > 3 {
		nY = 3
	} else if nY == 0 {
		nY = 1
	}

	endYPoint := startYPoint + rand.Intn(nY)

	c := b.randomColor()
	for i := startXPoint; i <= endXPoint; i++ {
		for j := startYPoint; j <= endYPoint; j++ {
			if show := rand.Intn(5); ((startYPoint-endYPoint) <= 1 && (endXPoint-startXPoint) <= 1) || show != 2 {
				b.pixels[j][i].SetColor(c)
				b.pixels[j][i].round = b.round
			}
		}
	}
}

func (b *Board) newRound() {
	b.round++
	b.shapeGenerator()
}

func (b *Board) Animate() {
	b.shapeGenerator()
	for range time.Tick(time.Millisecond * 250) {
		if b.isDone {
			return
		}
		if x := b.Down(); !x {
			b.checkForGameOver()
			b.checkForExplodeRow()
			b.newRound()
		}
	}
}

func (b *Board) Render() fyne.CanvasObject {
	items := make([]fyne.CanvasObject, 0)
	for i := 0; i < b.colPixel; i++ {
		b.pixels[i] = make([]*Pixel, b.rowPixel)
		for j := 0; j < b.rowPixel; j++ {
			b.pixels[i][j] = &Pixel{Rectangle: canvas.NewRectangle(color.White), round: 1}
			items = append(items, b.pixels[i][j].Rectangle)
		}
	}

	co := container.New(layout.NewGridLayout(b.rowPixel), items...)
	co.Resize(fyne.NewSize(float32(b.width), float32(b.height)))
	return co
}
