package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("GoTetris")

	const (
		width    = 500
		height   = 800
		perPixel = 50
	)

	board := NewBoard(width, height, perPixel)

	// --- UI: Score ---
	scoreLabel := canvas.NewText("Score: 0", color.Black)
	scoreLabel.TextSize = 24

	scoreBg := canvas.NewRectangle(color.White)
	scoreBg.Resize(fyne.NewSize(150, 40))
	scoreBg.Move(fyne.NewPos(20, 20))

	scoreLabel.Move(fyne.NewPos(30, 25))

	// --- Layout ---
	content := container.NewWithoutLayout(
		board.Render(),
		scoreBg,
		scoreLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(width, height))

	// --- Keyboard input ---
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		board.mu.Lock()
		defer board.mu.Unlock()

		if board.isDone {
			return
		}

		switch e.Name {
		case fyne.KeyLeft:
			board.Move(-1)

		case fyne.KeyRight:
			board.Move(1)

		case fyne.KeyDown:
			board.SoftDrop()

		case fyne.KeySpace:
			board.RotateCurrent()

		case fyne.KeyEscape:
			w.Close()
		}
	})

	// --- Result listener (score / game over) ---

	go func() {
		for r := range board.resultCh {
			switch r {
			case AddScoreResult:
				// score update
				board.mu.Lock()
				board.score++
				board.mu.Unlock()
				scoreLabel.Text = fmt.Sprintf("Score: %d", board.score)
				scoreLabel.Refresh()

			case GameOverResult:
				dialog.ShowInformation("Game Over", "Game over!!", w)
			}
		}
	}()

	// --- Start game loop ---
	go board.Animate()

	w.ShowAndRun()
}
