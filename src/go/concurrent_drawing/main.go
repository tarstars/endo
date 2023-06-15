package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"sync"
)

var (
	width  = 600
	height = 600
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Draw",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create a new image.RGBA value and fill it with colors
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			v := 0
			if x&y == 0 {
				v = 255
			}
			img.Set(x, y, color.RGBA{uint8(v), uint8(v), 0, 255})
		}
	}

	imageCh := make(chan *image.RGBA)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		x := 0
		y := 0
		meter := 0
		for {
			x++
			if x == width {
				x = 0
				y++
			}
			if y == height {
				x = 0
				y = 0
				break
			}

			img.Set(x, y, color.RGBA{255, 255, 255, 255})

			meter++
			if meter == 100 {
				imageCh <- img
				meter = 0
			}
		}
	}()

	go func() {
		wg.Wait()
		close(imageCh)
	}()

	for !win.Closed() {
		select {
		case img, ok := <-imageCh:
			if !ok {
				return
			}
			// Convert the RGBA image into a pixel Picture
			pic := pixel.PictureDataFromImage(img)

			// Create a sprite from the picture
			sprite := pixel.NewSprite(pic, pic.Bounds())

			// Clear the window to a white color
			win.Clear(colornames.White)

			// Draw the sprite to the window
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		default:
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
