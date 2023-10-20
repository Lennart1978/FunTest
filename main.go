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

var xMax, yMax float32

func checkWindowSize(w fyne.Window) {
	for {
		size := w.Canvas().Size()
		if xMax != size.Width || yMax != size.Height {
			xMax = size.Width
			yMax = size.Height
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	rand.NewSource(time.Now().UnixNano())
	myApp := app.New()
	w := myApp.NewWindow("Fun")
	w.Resize(fyne.NewSize(300, 600))
	w.SetFixedSize(false)
	w.CenterOnScreen()

	cont := container.NewWithoutLayout()

	go checkWindowSize(w)

	button := widget.NewButton("Start the fun !", func() {
		c := canvas.NewCircle(randomColor())
		c.StrokeWidth = 2
		c.Resize(fyne.NewSize(rand.Float32()*50.0, rand.Float32()*50.0))
		cont.Add(c)

		xStart, yStart := randomPosition()
		speed := randomSpeed()
		go spiralMotion(c, xStart, yStart, speed, w)
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

func spiralMotion(c *canvas.Circle, xStart, yStart, speed float64, w fyne.Window) {
	winSize := w.Canvas().Size()
	xMax = winSize.Width
	yMax = winSize.Height
	var theta float64
	var direction float64 = 1
	xMid, yMid := float64(xMax/2), float64(yMax/2)

	theta = math.Atan2(yStart-yMid, xStart-xMid)
	a := math.Sqrt((xStart-xMid)*(xStart-xMid)+(yStart-yMid)*(yStart-yMid)) - 10*theta

	circleSize := c.Size()
	radius := circleSize.Width / 2

	for {
		r := a + 10*theta
		newX := r*math.Cos(theta) + xMid
		newY := r*math.Sin(theta) + yMid

		if newX-float64(radius) <= 0 {
			newX = float64(radius)
			direction = -direction
		} else if newX+float64(radius) >= float64(xMax) {
			newX = float64(xMax) - float64(radius)
			direction = -direction
		}

		if newY-float64(radius) <= 0 {
			newY = float64(radius)
			direction = -direction
		} else if newY+float64(radius) >= float64(yMax) {
			newY = float64(yMax) - float64(radius)
			direction = -direction
		}

		c.Move(fyne.NewPos(float32(newX), float32(newY)))
		theta += speed * direction

		time.Sleep(time.Millisecond)
	}
}

func randomPosition() (float64, float64) {
	return rand.Float64() * float64(xMax-50), rand.Float64() * float64(yMax-100)
}

func randomSpeed() float64 {
	return rand.Float64()*0.01 + 0.005
}
