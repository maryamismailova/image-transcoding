package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

type ScalingImage struct {
	image  image.Image
	format string
}

func (sI *ScalingImage) Scale(ratio float64) (scaledImage *ScalingImage) {
	dst := image.NewRGBA(image.Rect(0, 0, int(float64(sI.image.Bounds().Max.X)*ratio), int(float64(sI.image.Bounds().Max.Y)*ratio)))

	draw.NearestNeighbor.Scale(dst, dst.Rect, sI.image, sI.image.Bounds(), draw.Over, nil)

	scaledImage = &ScalingImage{
		image:  dst,
		format: sI.format,
	}
	return scaledImage

}

func (sI *ScalingImage) Encode(w io.Writer) (err error) {
	switch sI.format {
	case "png":
		png.Encode(w, sI.image)
	case "jpeg":
		jpeg.Encode(w, sI.image, nil)
	default:
		return fmt.Errorf("format is unrecognized")
	}
	return nil
}

func NewImage(r io.Reader) (sI *ScalingImage, err error) {
	img, format, err := image.Decode(r)

	if err != nil {
		return nil, fmt.Errorf("unrecognized format")
	}

	sI = &ScalingImage{
		image:  img,
		format: format,
	}
	return sI, nil
}
