package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Tetris")

	const width int = 500
	const height int = 700
	const perPixel int = 50

	board := NewBoard(width, height, perPixel)
	board.SetWindowContent(w)

	w.Resize(fyne.NewSize(float32(width), float32(height)))

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyRight {
			board.Right()
		} else if event.Name == fyne.KeyLeft {
			board.Left()
		}
	})

	go board.Animate()

	w.ShowAndRun()
}
