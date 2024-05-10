package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Tetris")

	const width int = 500
	const height int = 800
	const perPixel int = 50

	board := NewBoard(width, height, perPixel)

	scoreLabel := widget.NewLabel("Score: 0")

	c := container.NewGridWithColumns(1,
		board.Render(),
		scoreLabel,
	)

	w.SetContent(c)

	w.Resize(fyne.NewSize(600, 900))

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
				scoreLabel.SetText(fmt.Sprintf("Score: %d", i))
				break
			}
		}
	}()

	go func() {
		board.Animate()
	}()

	w.ShowAndRun()
}
