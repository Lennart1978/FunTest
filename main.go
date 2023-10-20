package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	rand.NewSource(time.Now().UnixNano())
	myApp := app.New()
	w := myApp.NewWindow("Fun")
	w.Resize(fyne.NewSize(300, 600))
	w.SetFixedSize(true)

	cont := container.NewWithoutLayout()

	button := widget.NewButton("Start the fun !", func() {
		c := canvas.NewCircle(randomColor())
		c.StrokeWidth = 2
		c.Resize(fyne.NewSize(50, 50))
		cont.Add(c)

		xStart, yStart := randomPosition()
		speed := randomSpeed()
		go spiralMotion(c, xStart, yStart, speed)
	})

	cont.Add(button)
	button.Resize(fyne.NewSize(200, 50))

	w.SetContent(cont)
	w.ShowAndRun()
}

func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func spiralMotion(c *canvas.Circle, xStart, yStart, speed float64) {
	var theta float64
	var direction float64 = 1
	xMid, yMid := 150.0, 300.0

	theta = math.Atan2(yStart-yMid, xStart-xMid)
	a := math.Sqrt((xStart-xMid)*(xStart-xMid)+(yStart-yMid)*(yStart-yMid)) - 10*theta

	for {
		r := a + 10*theta
		x := r*math.Cos(theta) + xMid
		y := r*math.Sin(theta) + yMid
		c.Move(fyne.NewPos(float32(x), float32(y)))

		if x-25 <= 0 || x+25 >= 300 || y-25 <= 0 || y+25 >= 600 {
			direction = -direction
		}
		theta += speed * direction

		time.Sleep(time.Millisecond)
	}
}

func randomPosition() (float64, float64) {
	return rand.Float64()*200 + 50, rand.Float64()*400 + 100
}

func randomSpeed() float64 {
	return rand.Float64()*0.01 + 0.005
}
