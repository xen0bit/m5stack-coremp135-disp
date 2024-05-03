package main

import (
	"image"
	"image/draw"

	"github.com/gonutz/framebuffer"
)

const (
	screenWidth  = 320
	screenHeight = 240
	cols         = 320
	rows         = 240
	damping      = float32(0.95)
)

var (
	current  = [cols][rows]float32{}
	previous = [cols][rows]float32{}
)

func init() {
	for i := 0; i < cols-1; i++ {
		for j := 0; j < rows-1; j++ {
			current[i][j] = 0
			previous[i][j] = 0
		}
	}
	previous[screenHeight/2][screenWidth/2] = 255
}

func main() {
	fb, err := framebuffer.Open("/dev/fb1")
	if err != nil {
		panic(err)
	}
	defer fb.Close()

	noiseImage := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	//Render loop
	for {
		//Kernel Math for simple water animation
		for i := 1; i < cols-1; i++ {
			for j := 1; j < rows-1; j++ {
				current[i][j] = (previous[i-1][j]+previous[i+1][j]+
					previous[i][j-1]+previous[i][j+1]+
					previous[i-1][j-1]+previous[i-1][j+1]+
					previous[i+1][j-1]+previous[i+1][j+1])/4 - current[i][j]

				current[i][j] = current[i][j] * damping
				index := (i + j*cols) * 4

				noiseImage.Pix[index+0] = 0
				noiseImage.Pix[index+1] = 0
				noiseImage.Pix[index+2] = uint8(current[i][j] * 255)
				noiseImage.Pix[index+3] = 255
			}
		}
		//swap
		temp := previous
		previous = current
		current = temp
		//draw
		draw.Draw(fb, fb.Bounds(), noiseImage, image.ZP, draw.Src)
	}
}
