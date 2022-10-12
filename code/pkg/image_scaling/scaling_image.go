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

func (sI *ScalingImage) GetFormat() string {
	return sI.format
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
func ScaleImageFromSource(sourcePath string, destPath string, scaleY int, scaleX int) (*ScalingImage, error) {
	fSrc, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("%v: failed opening source file %s", err, sourcePath)
	}
	defer fSrc.Close()

	fDest, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to create destination file %s", err, destPath)
	}
	defer fDest.Close()

	scaledImage, err := ScaleImage(fSrc, fDest, scaleY, scaleX)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to scale the image %s", err, sourcePath)
	}

	log.Printf("Finished image transcoding")
	return scaledImage, nil
}

func ScaleImage(reader io.Reader, writer io.Writer, scaleY int, scaleX int) (*ScalingImage, error) {

	img, err := NewImage(reader)
	if err != nil {
		return nil, fmt.Errorf("%v: error reading image bytes", err)
	}
	log.Printf("Image of format: %s\n", img.format)

	squareResolution := int(math.Min(float64(scaleX), float64(scaleY)))
	scaleRatio := math.Min(float64(scaleX)/float64(img.image.Bounds().Max.X), float64(scaleY)/float64(img.image.Bounds().Max.Y))
	log.Printf("The square's resolution: %d:%d", squareResolution, squareResolution)
	log.Printf("Scale ratio to be applied is %f\n", scaleRatio)

	scaledImage := img.Scale(scaleRatio)

	err = scaledImage.Encode(writer)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to encode image", err)
	}

	return scaledImage, nil
}
