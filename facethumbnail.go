package facethumbnail

// TODO: Add the thumbnail logic
import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	//"github.com/nfnt/resize"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// ResizeImage uses an instance of FaceDetector to detect face in srcPath and generates a thumbnail of size x size in dstPath
func ResizeImage(fd *FaceDetector, srcPath, dstPath string, size uint) error {
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

	log.Printf("Opened image %v or size (%v,%v), cropping to (%v,%v)", srcPath, img.Bounds().Dx(), img.Bounds().Dy(), size, size)

	// default center is mid-point of image
	var faceCenter image.Point = image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)

	// Now if a face detector is provided and if that detector detects a face use
	// the center of the detected face to center the thumbnail image
	if fd != nil {
		faces, err := fd.DetectFacesInImageFile(srcPath)
		if err != nil {
			return fmt.Errorf("Face detection failed with %v", err)
		}

		nFaces := len(faces)
		log.Printf("Detect %v faces", nFaces)

		largestFaceSize := 0

		for _, face := range faces {
			log.Printf("Detected face %v", face)
			faceSize := face.Dx() * face.Dy()
			if faceSize > largestFaceSize {
				largestFaceSize = faceSize
				x := (face.Min.X + face.Max.X) / 2
				y := (face.Min.Y + face.Max.Y) / 2

				// center of detected face
				faceCenter = image.Pt(x, y)
			}
		}
	}

	log.Printf("Using faceCenter %v", faceCenter)

	// In the code below we are attempting to find a square whose center is close to the center of the found face

	// Generate a square image of sizeSquare which is equal to the width or height of the image
	// whichever is smaller
	sizeSquare := min(img.Bounds().Dx(), img.Bounds().Dy())

	// X, Y are the left top corner of cropped image of size sizeSquare

	// Attempt to use x so that facecenter.X is exactly at the center of the sizeSquare image
	x := faceCenter.X - (sizeSquare / 2)
	// now if x is negative then the center was far to the left, so use x so that we take all of the
	// image from the very left end
	if x < 0 {
		x = 0
	} else if (x + sizeSquare) > img.Bounds().Dx() {
		// if x + sizeSquare is beyond the bounds then we need to move x more left so that we can
		// go till the end of the image and no more
		x = img.Bounds().Dx() - sizeSquare
	}

	// same logic as X, but this time for vertical (Y) axis
	y := faceCenter.Y - (sizeSquare / 2)
	if y < 0 {
		y = 0
	} else if (y + sizeSquare) > img.Bounds().Dy() {
		y = img.Bounds().Dy() - sizeSquare
	}

	// Now crop the image with anchor which is the left-top coordinate of the cropped image
	anchor := image.Pt(x, y)
	config := cutter.Config{
		Width:  sizeSquare,
		Height: sizeSquare,
		Anchor: anchor,
		Mode:   cutter.TopLeft, // optional, default value
	}

	croppedImg, err := cutter.Crop(img, config)

	log.Printf("Log config %+v", config)

	resizedImage := resize.Resize(size, size, croppedImg, resize.Lanczos3)

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, resizedImage, nil)
	log.Printf("Generated %v", dstPath)

	return nil
}

func min(a, b int) (r int) {
	if a < b {
		return a
	}
	return b
}
