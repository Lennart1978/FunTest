package main

import (
	"encoding/hex"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"bytes"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var xMax, yMax float32
var bt *widget.Button

type TappableCircle struct {
	widget.BaseWidget
	Color color.Color
	Cont  *fyne.Container
}

func NewTappableCircle(col color.Color, size fyne.Size, cont *fyne.Container) *TappableCircle {
	c := &TappableCircle{Color: col, Cont: cont}
	c.ExtendBaseWidget(c)
	return c
}

func (c *TappableCircle) CreateRenderer() fyne.WidgetRenderer {
	circ := canvas.NewCircle(c.Color)
	// Größe des Kreises direkt setzen
	size := c.BaseWidget.Size()
	if size.IsZero() {
		// Verwenden Sie eine Standardgröße, wenn die Größe nicht gesetzt wurde
		size = fyne.NewSize(50, 50)
	}
	circ.Resize(size)
	return &circleRenderer{circle: circ, objects: []fyne.CanvasObject{circ}, size: size} // Größe auch hier speichern
}

type circleRenderer struct {
	circle  *canvas.Circle
	objects []fyne.CanvasObject
	size    fyne.Size // Speichern Sie die Größe hier
}

func (c *circleRenderer) Layout(size fyne.Size) {
	c.circle.Resize(size)
}

func (c *circleRenderer) MinSize() fyne.Size {
	return c.size
}

func (c *circleRenderer) Refresh() {
	c.circle.FillColor = c.circle.FillColor
	canvas.Refresh(c.circle)
}

func (c *circleRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (c *circleRenderer) Objects() []fyne.CanvasObject {
	return c.objects
}

func (c *circleRenderer) Destroy() {}

func (c *TappableCircle) Tapped(*fyne.PointEvent) {
	c.Cont.Remove(c)
}

func (c *TappableCircle) TappedSecondary(*fyne.PointEvent) {}

func checkWindowSize(w fyne.Window, c fyne.CanvasObject) {
	for {
		size := w.Canvas().Size()
		if xMax != size.Width || yMax != size.Height {
			xMax = size.Width
			yMax = size.Height
			c.Resize(fyne.NewSize(xMax, yMax))
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	mp3Data, err := hex.DecodeString(mp3DataHex)
	if err != nil {
		panic(err)
	}

	go playMP3(mp3Data)
	rand.NewSource(time.Now().UnixNano())
	myApp := app.New()
	w := myApp.NewWindow("Fun")
	//w.Resize(fyne.NewSize(300, 600))
	w.SetFullScreen(true)
	w.SetFixedSize(false)

	cont := container.NewWithoutLayout()

	gradient := canvas.NewRadialGradient(color.Black, randomColor())
	cont.Add(gradient)
	gradient.Resize(fyne.NewSize(300, 600))

	go checkWindowSize(w, gradient)

	button := widget.NewButton("Spiral !", func() {
		c := NewTappableCircle(randomColor(), fyne.NewSize(rand.Float32()*50.0, rand.Float32()*50.0), cont)

		c.Resize(fyne.NewSize(rand.Float32()*50.0, rand.Float32()*50.0))
		cont.Add(c)

		xStart, yStart := randomPosition()
		speed := randomSpeed()
		go spiralMotion(c, xStart, yStart, speed, w)
	})

	buttonAngular := widget.NewButton("Angular !", func() {
		c := NewTappableCircle(randomColor(), fyne.NewSize(rand.Float32()*50.0, rand.Float32()*50.0), cont)

		c.Resize(fyne.NewSize(rand.Float32()*50.0, rand.Float32()*50.0))
		cont.Add(c)

		xStart, yStart := randomPosition()
		speed := randomSpeed()
		go angularMotion(c, xStart, yStart, speed, w)
	})

	buttonFC := widget.NewButton("Full Screen", func() {
		if w.FullScreen() {
			w.SetFullScreen(false)
		} else {
			w.SetFullScreen(true)
		}
	})
	cont.Add(buttonFC)

	buttonClear := widget.NewButton("Kill !", func() {
		if (len(cont.Objects)) == 4 {
			bt.Move(fyne.NewPos(xMax/2-50, yMax/2-25))
			gameOver(bt)
		}
		if (len(cont.Objects)) == 3 {
			myApp.Quit()
			os.Exit(0)
		}
		cont.Remove(cont.Objects[len(cont.Objects)-1])

	})

	bt = buttonClear
	cont.Add(buttonClear)

	cont.Add(button)
	button.Resize(fyne.NewSize(150, 50))
	cont.Add(buttonAngular)
	buttonClear.Resize(fyne.NewSize(100, 50))
	buttonClear.Move(fyne.NewPos(155, 0))
	buttonFC.Resize(fyne.NewSize(100, 50))
	buttonFC.Move(fyne.NewPos(260, 0))
	buttonAngular.Resize(fyne.NewSize(100, 50))
	buttonAngular.Move(fyne.NewPos(365, 0))
	w.SetContent(cont)
	w.ShowAndRun()
}
func gameOver(b *widget.Button) {
	b.SetText("Game Over !")
}

func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func angularMotion(c *TappableCircle, xStart, yStart, speed float64, w fyne.Window) {
	winSize := w.Canvas().Size()
	xMax, yMax := float64(winSize.Width), float64(winSize.Height)
	speed *= 30.0
	var xSpeed, ySpeed float64 = speed, speed
	circleSize := c.Size()
	radius := float64(circleSize.Width / 2)

	rand.NewSource(time.Now().UnixNano())

	for {
		xStart += xSpeed
		yStart += ySpeed

		if xStart-radius <= 0 {
			xStart = radius
			xSpeed = speed * (1 + rand.Float64()) // Zufällige Geschwindigkeitskomponente in X-Richtung
		} else if xStart+radius >= xMax {
			xStart = xMax - radius
			xSpeed = -speed * (1 + rand.Float64()) // Zufällige Geschwindigkeitskomponente in X-Richtung
		}

		if yStart-radius <= 0 {
			yStart = radius
			ySpeed = speed * (1 + rand.Float64()) // Zufällige Geschwindigkeitskomponente in Y-Richtung
		} else if yStart+radius >= yMax {
			yStart = yMax - radius
			ySpeed = -speed * (1 + rand.Float64()) // Zufällige Geschwindigkeitskomponente in Y-Richtung
		}

		c.Move(fyne.NewPos(float32(xStart), float32(yStart)))

		time.Sleep(time.Millisecond)
	}
}

func spiralMotion(c *TappableCircle, xStart, yStart, speed float64, w fyne.Window) {
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

	rebound := 5.0
	for {
		r := a + 10*theta
		newX := r*math.Cos(theta) + xMid
		newY := r*math.Sin(theta) + yMid

		if newX-float64(radius) <= 0 {
			newX = float64(radius) + rebound
			direction = -direction
		} else if newX+float64(radius) >= float64(xMax) {
			newX = float64(xMax) - float64(radius) - rebound
			direction = -direction
		}

		if newY-float64(radius) <= 0 {
			newY = float64(radius) + rebound
			direction = -direction
		} else if newY+float64(radius) >= float64(yMax) {
			newY = float64(yMax) - float64(radius) - rebound
			direction = -direction
		}

		c.Move(fyne.NewPos(float32(newX), float32(newY)))
		theta += speed * direction

		time.Sleep(time.Millisecond)
	}
}

func randomPosition() (float64, float64) {
	return rand.Float64() * float64(xMax), rand.Float64() * float64(yMax)
}

func randomSpeed() float64 {
	return rand.Float64()*0.01 + 0.005
}

func playMP3(data []byte) {
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to create decoder: %v", err)
	}

	p, err := oto.NewContext((int(decoder.SampleRate())), 2, 2, 8192)
	if err != nil {
		log.Fatalf("failed to create player: %v", err)
	}
	defer p.Close()

	if _, err := copyBuffer(p.NewPlayer(), decoder); err != nil {
		log.Fatalf("failed to copy buffer: %v", err)
	}

}

func copyBuffer(dst *oto.Player, src *mp3.Decoder) (int64, error) {
	buf := make([]byte, 8192)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
			written += int64(nw)
		}
		if er != nil {
			if er == io.EOF {
				return written, nil
			}
			return written, er
		}
	}
}
