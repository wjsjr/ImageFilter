// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The Image represents a structure for working with PNG images.
// From Professor Samuels: You are allowed to update this and change it as you wish!
type Image struct {
	in     *image.RGBA64   //The original pixels before applying the effect
	out    *image.RGBA64   //The updated pixels after applying teh effect
	Bounds image.Rectangle //The size of the image
}

//
// Public functions
//

//Helper function which moves img.in pointer to point to img.out and makes img.out an empty image
func (img *Image) SwapIn() {
	img.in = img.out
	img.out = image.NewRGBA64(img.in.Bounds())
}

//Helper function which moves img.out pointer to point to img.in
func (img *Image) SwapOut() {
	img.out = img.in
}

// Load returns a Image that was loaded based on the filePath parameter
// From Professor Samuels:  You are allowed to modify and update this as you wish
func Load(filePath string) (*Image, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inOrig, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	bounds := inOrig.Bounds()

	outImg := image.NewRGBA64(bounds)
	inImg := image.NewRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inOrig.At(x, y).RGBA()
			inImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	task := &Image{}
	task.in = inImg
	task.out = outImg
	task.Bounds = bounds
	return task, nil
}

// Save saves the image to the given file
// From Professor Samuels:  You are allowed to modify and update this as you wish
func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
	if err != nil {
		return err
	}
	return nil
}

//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

//Helper function which returns height of an image
func (img *Image) Height() int {
	return img.Bounds.Dy()
}
