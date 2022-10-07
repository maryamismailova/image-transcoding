package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

func scaleImageFromSource(sourcePath string, destPath string, scaleY int, scaleX int) error {
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

	scaleRatio := math.Min(float64(scaleX)/float64(img.image.Bounds().Max.X), float64(scaleY)/float64(img.image.Bounds().Max.Y))
	log.Printf("Scale ratio to be applied is %f\n", scaleRatio)

	scaledImage := img.Scale(scaleRatio)

	fDest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fDest.Close()
	err = scaledImage.Encode(fDest)
	if err != nil {
		return fmt.Errorf("%v: failed to create destination file %s", err, destPath)
	}

	return nil
}

func main() {
	var err error
	err = scaleImageFromSource("tests/test-1.png", "outputs/test-1.png", 400, 800)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	err = scaleImageFromSource("tests/test-1.gif", "outputs/test-1.gif", 400, 800)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	err = scaleImageFromSource("tests/test-1.jpg", "outputs/test-1.jpg", 400, 800)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

}
