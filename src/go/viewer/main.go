package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/tarstars/endo/src/go/aide"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Rna viewer",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	aide.Must(err)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	var testText *text.Text
	testText = text.New(pixel.V(10, 100), atlas)
	testText.Color = colornames.Black
	_, err = fmt.Fprintln(testText, "Good news everyone!")
	aide.Must(err)

	imd := imdraw.New(nil)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		testText.Draw(win, pixel.IM.Scaled(testText.Orig, 2))

		for r := 0.; r < 101; r += 2 {
			imd.Color = colornames.Yellow
			imd.Push(pixel.V(win.Bounds().Max.X/2, win.Bounds().Max.Y/2))
			imd.Circle(r, 1)
			imd.Draw(win)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
