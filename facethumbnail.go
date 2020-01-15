package facethumbnail

// TODO: Add the thumbnail logic
import (
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func ResizeImage(srcPath, dstPath string, size uint) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	file.Close()
	w := uint(0)
	h := uint(0)
	if img.Bounds().Dx() > img.Bounds().Dy() {
		h = size
		w = 0
	} else {
		w = size
		h = 0
	}

	// TODO: Support portrait thumbnail that uses top part of photo
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio

	resizedImage := resize.Resize(w, h, img, resize.Lanczos3)
	croppedImg, err := cutter.Crop(resizedImage, cutter.Config{
		Width:  int(size),
		Height: int(size),
		//Anchor: image.Point{100, 100},
		Mode: cutter.Centered, // optional, default value
	})

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, croppedImg, nil)
	log.Printf("Generated %v", dstPath)

	return nil
}
