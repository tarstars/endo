package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	_ "image/png"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Draw",
		Bounds: pixel.R(0, 0, 600, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create a new image.RGBA value and fill it with colors
	img := image.NewRGBA(image.Rect(0, 0, 600, 600))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			v := 0
			if x&y == 0 {
				v = 255
			}
			img.Set(x, y, color.RGBA{uint8(v), uint8(v), 0, 255})
		}
	}

	x := 0
	y := 0

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)

		for meter := 70; meter > 0; meter-- {
			x++
			if x == 600 {
				x = 0
				y++
			}
			if y == 600 {
				x = 0
				y = 0
			}

			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}

		// Convert the RGBA image into a pixel Picture
		pic := pixel.PictureDataFromImage(img)

		// Create a sprite from the picture
		sprite := pixel.NewSprite(pic, pic.Bounds())

		// Draw the sprite to the window
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
