package main

import (
	"image/color"
	"encoding/hex"
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
	mp3Data, err := hex.DecodeString(mp3DataHex)
     if err != nil {
         panic(err)
     }
	
	go playMP3(mp3Data)
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

	buttonClear := widget.NewButton("Kill !", func() {
		if (len(cont.Objects)) == 2 {
			bt.Move(fyne.NewPos(xMax/2-100, yMax/2-50))
			gameOver(bt)
		}
		if (len(cont.Objects)) == 1 {
			myApp.Quit()
			os.Exit(0)
		}
		cont.Remove(cont.Objects[len(cont.Objects)-1])

	})
	bt = buttonClear
	cont.Add(buttonClear)

	cont.Add(button)
	button.Resize(fyne.NewSize(150, 50))
	buttonClear.Resize(fyne.NewSize(100, 50))
	buttonClear.Move(fyne.NewPos(155, 0))
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
	// MP3-Decoder erstellen
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to create decoder: %v", err)
	}

	// Oto-Player-Context erstellen
	p, err := oto.NewContext((int(decoder.SampleRate())), 2, 2, 8192)
	if err != nil {
		log.Fatalf("failed to create player: %v", err)
	}

	// MP3-Datei abspielen
	if _, err := copyBuffer(p.NewPlayer(), decoder); err != nil {
		log.Fatalf("failed to copy buffer: %v", err)
	}
}

// copyBuffer kopiert Daten vom Decoder zum Player
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
			return written, er
		}
	}
}
