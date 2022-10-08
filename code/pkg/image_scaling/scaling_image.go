// Package image_scaling provides a code base for transcoding images
package image_scaling

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"
)

// ScalingImage is a placeholder for image.Image as well as its format
// created to avoid constant decoding of image
// as well as to create functions linked to it
type ScalingImage struct {
	image  image.Image
	format string
}

// Scale function given a ration scales the source image
// and returns a new ScalingImage object
func (sI *ScalingImage) Scale(ratio float64) (scaledImage *ScalingImage) {
	dst := image.NewRGBA(image.Rect(0, 0, int(float64(sI.image.Bounds().Max.X)*ratio), int(float64(sI.image.Bounds().Max.Y)*ratio)))

	draw.NearestNeighbor.Scale(dst, dst.Rect, sI.image, sI.image.Bounds(), draw.Over, nil)

	scaledImage = &ScalingImage{
		image:  dst,
		format: sI.format,
	}
	return scaledImage

}

// Encode the bytes into an image
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

// NewImage function reads image object from bytes stream
// This function also errors in case unaccepted image format is provided to it!
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

// ScaleImageFromSource scales source image
// given the source and destination paths, as well as new resolution
func ScaleImageFromSource(sourcePath string, destPath string, scaleY int, scaleX int) error {
	log.Printf("Start scaling of %s to %s with expected resolution %d:%d \n", sourcePath, destPath, scaleY, scaleX)
	fSrc, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("%v: failed opening source file %s", err, sourcePath)
	}
	defer fSrc.Close()

	img, err := NewImage(fSrc)
	if err != nil {
		return fmt.Errorf("%v: error reading %s", err, sourcePath)
	}
	log.Printf("%s image of format: %s\n", sourcePath, img.format)

	squareResolution := int(math.Min(float64(scaleX), float64(scaleY)))
	log.Printf("The square's resolution: %d:%d", squareResolution, squareResolution)
	scaleRatio := math.Min(float64(scaleX)/float64(img.image.Bounds().Max.X), float64(scaleY)/float64(img.image.Bounds().Max.Y))
	log.Printf("Scale ratio to be applied is %f\n", scaleRatio)

	scaledImage := img.Scale(scaleRatio)

	fDest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("%v: failed to create destination file %s", err, destPath)
	}
	defer fDest.Close()
	err = scaledImage.Encode(fDest)
	if err != nil {
		return fmt.Errorf("%v: failed to encode to %s", err, destPath)
	}
	log.Printf("Finished image transcoding")
	return nil
}
