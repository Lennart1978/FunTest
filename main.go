package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Fun")
	w.Resize(fyne.NewSize(300, 600))
	w.SetFixedSize(true)
	c := canvas.NewCircle(color.RGBA{255, 0, 0, 255})
	c.StrokeWidth = 2
	c.Resize(fyne.NewSize(50, 50))

	type direct struct {
		left, right, up, down bool
	}
	var direction direct

	button := widget.NewButton("Start the fun !", func() {
		go func() {
			var x, y float32
			direction.left = false
			direction.right = true
			direction.up = false
			direction.down = true

			x = 150.0
			y = 300.0
			for {
				c.Move(fyne.NewPos(x, y))
				if direction.left {
					x = x - 1.0
				}
				if direction.right {
					x = x + 1.0
				}
				if direction.up {
					y = y - 1.0
				}
				if direction.down {
					y = y + 1.0
				}
				if x < 0 {
					direction.left = false
					direction.right = true
				}
				if x > 300 {
					direction.right = false
					direction.left = true
				}
				if y <= 0 {
					direction.up = false
					direction.down = true
				}
				if y > 600 {
					direction.down = false
					direction.up = true
				}
				time.Sleep(time.Millisecond)
			}
		}()
	})
	cont := container.NewWithoutLayout(c, button)
	button.Resize(fyne.NewSize(200, 50))
	c.Move(fyne.NewPos(100, 100))

	w.SetContent(cont)

	w.ShowAndRun()
}
