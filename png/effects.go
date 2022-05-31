// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"fmt"
	"image/color"
	"os"
)

//Helper function which takes an image and effect and applies effect to image
func ApplyEffect(image *Image, e string) {

	//If effect is grayscale, call grayscale method. Else call Convolute method with correct Kernel
	if e == "G" {
		image.Grayscale()
	} else if e == "S" {
		image.Convolute([][]float64{{0, -1, 0}, {-1, 5, -1}, {0, -1, 0}})
	} else if e == "E" {
		image.Convolute([][]float64{{-1, -1, -1}, {-1, 8, -1}, {-1, -1, -1}})
	} else if e == "B" {
		image.Convolute([][]float64{{1 / 9.0, 1 / 9.0, 1 / 9.0}, {1 / 9.0, 1 / 9.0, 1 / 9.0}, {1 / 9.0, 1 / 9.0, 1 / 9.0}})
	} else {
		fmt.Println("ERROR: INVALID effect: ", e)
		os.Exit(1)
	}
}

//Function which applies 2d convolution: https://en.wikipedia.org/wiki/Kernel_(image_processing)
func (img *Image) Convolute(k [][]float64) {
	bounds := img.out.Bounds()
	var r, g, b, a uint32

	//Iterate over each pixel in image.in
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum float64

			//Apply kernel to pixel, keeping track of sum of kernel impact on the pixel
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					cy := y + ky - 1
					cx := x + kx - 1

					//If cx/cy out of bounds, apply zero padding
					if cx >= 0 && cy >= 0 {
						r, g, b, a = img.in.At(cx, cy).RGBA()
					} else {
						r = 0
						g = 0
						b = 0
					}
					rSum += float64(r) * k[ky][kx]
					gSum += float64(g) * k[ky][kx]
					bSum += float64(b) * k[ky][kx]
				}
			}
			//Write to img.out
			img.out.Set(x, y, color.RGBA64{clamp(rSum), clamp(gSum), clamp(bSum), uint16(a)})
		}
	}

}

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale() {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}
