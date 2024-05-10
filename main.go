package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"image/color"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Tetris")

	const width int = 500
	const height int = 800
	const perPixel int = 50

	board := NewBoard(width, height, perPixel)

	scoreLabel := canvas.NewText("Score: 0", color.Black)
	scoreLabel.TextSize = 32
	scoreContainer := canvas.NewRectangle(color.White)
	scoreContainer.Resize(fyne.NewSize(150, 50))
	scoreContainer.FillColor = color.White
	scoreContainer.Move(fyne.NewPos(25, 25))
	scoreLabel.Move(fyne.NewPos(35, 25))
	scoreLabel.Resize(fyne.NewSize(scoreLabel.Size().Height+10, scoreLabel.Size().Height+10))

	c := container.NewWithoutLayout(
		board.Render(),
		scoreContainer,
		scoreLabel,
	)
	w.SetContent(c)

	w.Resize(fyne.NewSize(float32(width), float32(height)))

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyRight {
			board.Right()
		} else if event.Name == fyne.KeyLeft {
			board.Left()
		}
	})

	i := 0

	go func() {
		for {
			switch <-board.resultCh {
			case GameOverResult:
				dialog.ShowInformation("Game Over", "Game over !!", w)
				break
			case AddScoreResult:
				i++
				scoreLabel.Text = fmt.Sprintf("Score: %d", i)
				scoreLabel.Refresh()
				break
			}
		}
	}()

	go func() {
		board.Animate()
	}()

	w.ShowAndRun()
}
