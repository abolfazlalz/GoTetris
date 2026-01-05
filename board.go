package main

import (
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	rowPixel int // width
	colPixel int // height

	mu sync.Mutex

	pixels  [][]*Pixel      // render layer
	fixed   [][]color.Color // locked blocks
	current *Tetromino      // active piece

	isDone   bool
	resultCh chan Result
	score    int
}

func NewBoard(width, height, perPixel int) *Board {
	rows := height / perPixel
	cols := width / perPixel

	b := &Board{
		width:    width,
		height:   height,
		perPixel: perPixel,
		rowPixel: cols,
		colPixel: rows,
		resultCh: make(chan Result, 1),
	}

	b.fixed = make([][]color.Color, rows)
	for y := 0; y < rows; y++ {
		b.fixed[y] = make([]color.Color, cols)
		for x := 0; x < cols; x++ {
			b.fixed[y][x] = color.White
		}
	}

	return b
}

func (b *Board) Animate() {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		if b.isDone {
			b.mu.Unlock()
			return
		}
		b.tick()
		b.mu.Unlock()
	}
}

func (b *Board) tick() {
	if b.current == nil {
		b.spawnTetromino()
		return
	}

	if b.canPlace(b.current, b.current.Shape, b.current.X, b.current.Y+1) {
		b.current.Y++
	} else {
		b.lockCurrent()
	}

	b.redraw()
}

func (b *Board) spawnTetromino() {
	t := NewRandomTetromino(b.rowPixel)
	b.current = t

	if !b.canPlace(t, t.Shape, t.X, t.Y) {
		b.isDone = true
		b.resultCh <- GameOverResult
	}
}

func (b *Board) lockCurrent() {
	t := b.current

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if t.Shape[y][x] == 1 {
				b.fixed[t.Y+y][t.X+x] = t.Color
			}
		}
	}

	b.clearLines()
	b.current = nil
}

func (b *Board) clearLines() {
	for y := b.colPixel - 1; y >= 0; y-- {
		full := true
		for x := 0; x < b.rowPixel; x++ {
			if b.fixed[y][x] == color.White {
				full = false
				break
			}
		}

		if full {
			b.score++
			b.resultCh <- AddScoreResult

			for yy := y; yy > 0; yy-- {
				b.fixed[yy] = append([]color.Color{}, b.fixed[yy-1]...)
			}
			for x := 0; x < b.rowPixel; x++ {
				b.fixed[0][x] = color.White
			}
			y++
		}
	}
}

func (b *Board) RotateCurrent() {
	if b.current == nil {
		return
	}

	rotated := rotateMatrixCW(b.current.Shape)
	if b.canPlace(b.current, rotated, b.current.X, b.current.Y) {
		b.current.Shape = rotated
		b.redraw()
	}
}

func (b *Board) Move(dx int) {
	if b.current == nil {
		return
	}

	if b.canPlace(b.current, b.current.Shape, b.current.X+dx, b.current.Y) {
		b.current.X += dx
		b.redraw()
	}
}

func (b *Board) SoftDrop() {
	if b.current == nil {
		return
	}

	if b.canPlace(b.current, b.current.Shape, b.current.X, b.current.Y+1) {
		b.current.Y++
		b.redraw()
	}
}

func (b *Board) canPlace(t *Tetromino, shape [][]int, nx, ny int) bool {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if shape[y][x] == 0 {
				continue
			}

			bx := nx + x
			by := ny + y

			if bx < 0 || bx >= b.rowPixel || by < 0 || by >= b.colPixel {
				return false
			}

			if b.fixed[by][bx] != color.White {
				return false
			}
		}
	}
	return true
}

func (b *Board) redraw() {
	for y := 0; y < b.colPixel; y++ {
		for x := 0; x < b.rowPixel; x++ {
			b.pixels[y][x].SetColor(b.fixed[y][x])
		}
	}

	if b.current == nil {
		return
	}

	t := b.current
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if t.Shape[y][x] == 1 {
				b.pixels[t.Y+y][t.X+x].SetColor(t.Color)
			}
		}
	}
}

func (b *Board) Render() fyne.CanvasObject {
	items := make([]fyne.CanvasObject, 0, b.rowPixel*b.colPixel)

	b.pixels = make([][]*Pixel, b.colPixel)
	for y := 0; y < b.colPixel; y++ {
		b.pixels[y] = make([]*Pixel, b.rowPixel)
		for x := 0; x < b.rowPixel; x++ {
			p := &Pixel{
				Rectangle: canvas.NewRectangle(color.White),
			}
			b.pixels[y][x] = p
			items = append(items, p.Rectangle)
		}
	}

	grid := container.NewGridWithColumns(b.rowPixel, items...)
	grid.Resize(fyne.NewSize(float32(b.width), float32(b.height)))

	return grid
}
