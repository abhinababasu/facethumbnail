package facethumbnail

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// faceDetector interface defines face detection interface ResizeImage expects
type FaceDetector interface {
	// Init initializes the face detector with maximum and minimum size of expected face
	Init(minSize, maxSize int) error

	// DetectFacesInImageFile detects a list of faces in an image file
	DetectFacesInImageFile(sourceFilePath string) ([]image.Rectangle, error)
}

// ResizeResult contains details of the resizing done
type ResizeResult struct {
	Center     image.Point
	FacesCount int
}

// ResizeToAspectRatio crops image to meet the aspect ratio of w:h
func ResizeToAspectRatio(fd FaceDetector, srcPath, dstPath string, w, h uint) (ResizeResult, error) {
	var result ResizeResult
	file, err := os.Open(srcPath)
	if err != nil {
		return result, err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return result, err
	}

	file.Close()

	log.Printf("Opened image %v of size (%v,%v), cropping to (%v,%v)", srcPath, img.Bounds().Dx(), img.Bounds().Dy(), w, h)

	faceCenter, nFaces, err := getFaceCenter(fd, srcPath, img)
	if err != nil {
		return result, err
	}

	result.FacesCount = nFaces

	// At this point faceCenter is either the center of the source image or the center of the
	// largest face detected in the image
	log.Printf("Using faceCenter %v", faceCenter)
	result.Center = faceCenter

	// In the code below we are attempting to find a square whose center is close to the center of the detected face

	// Generate a square image of sizeSquare which is equal to the width or height of the image
	// whichever is smaller
	imgwf, imghf := float32(img.Bounds().Dx()), float32(img.Bounds().Dy())

	// Find new dimentions after crop
	ratio := float32(w) / float32(h)
	imghfinal := imghf
	imgwfinal := imghfinal * ratio

	// calcuated width being larger than source width means that the image is too narrow for the crop
	// so go reverse, use the width and crop height
	if imgwfinal > imgwf {
		imgwfinal = imgwf
		imghfinal = imgwfinal / ratio
	}

	wfinal, hfinal := int(imgwfinal), int(imghfinal)

	err = crop(img, dstPath, faceCenter, wfinal, hfinal)

	return result, err
}

// ResizeImage uses an instance of FaceDetector to detect face in srcPath and generates a thumbnail of sizexsize in dstPath
// If no facedetector is given or no faces are detected, then the center of image is used for the thumbnail
func ResizeImage(fd FaceDetector, srcPath, dstPath string, size uint) (ResizeResult, error) {
	var result ResizeResult
	file, err := os.Open(srcPath)
	if err != nil {
		return result, err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return result, err
	}

	file.Close()

	log.Printf("Opened image %v of size (%v,%v), cropping to (%v,%v)", srcPath, img.Bounds().Dx(), img.Bounds().Dy(), size, size)

	faceCenter, nFaces, err := getFaceCenter(fd, srcPath, img)
	if err != nil {
		return result, err
	}

	result.FacesCount = nFaces

	// At this point faceCenter is either the center of the source image or the center of the
	// largest face detected in the image
	log.Printf("Using faceCenter %v", faceCenter)
	result.Center = faceCenter

	// In the code below we are attempting to find a square whose center is close to the center of the detected face

	// Generate a square image of sizeSquare which is equal to the width or height of the image
	// whichever is smaller
	sizeSquare := min(img.Bounds().Dx(), img.Bounds().Dy())
	err = crop(img, dstPath, faceCenter, sizeSquare, sizeSquare)

	return result, err
}

func min(a, b int) (r int) {
	if a < b {
		return a
	}
	return b
}

// getFaceCenter finds the rectangle for the face that is to be kept in the cropped picture
func getFaceCenter(fd FaceDetector, srcPath string, img image.Image) (image.Point, int, error) {
	// default center is mid-point of image
	var faceCenter image.Point = image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
	nFaces := 0
	// Now if a face detector is provided and if that detector detects a face use
	// the center of the detected face to center the thumbnail image
	if fd != nil {
		faces, err := fd.DetectFacesInImageFile(srcPath)
		if err != nil {
			return faceCenter, -1, fmt.Errorf("Face detection failed with %v", err)
		}

		nFaces = len(faces)
		log.Printf("Detected %v faces", nFaces)

		largestFaceSize := 0

		// Iterate through the faces and find the largest detected face and use that
		//
		// NOTE: Tried other mechanisms like union of all detected faces and taking the center of it
		// but that did not yield good results (e.g. picked empty space in-between two persons in a portrait)
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

	return faceCenter, nFaces, nil
}

// crop image into final size while ensuring the face is in the picture
func crop(img image.Image, dstPath string, faceCenter image.Point, wfinal, hfinal int) error {

	// X, Y are the left top corner of cropped image of size sizeSquare
	// Attempt to use x so that facecenter.X is exactly at the center of the sizeSquare image
	x := faceCenter.X - (wfinal / 2)
	// now if x is negative then the center was far to the left, so use x so that we take all of the
	// image from the very left end
	if x < 0 {
		x = 0
	} else if (x + wfinal) > img.Bounds().Dx() {
		// if x + sizeSquare is beyond the bounds then we need to move x more left so that we can
		// go till the end of the image and no more
		x = img.Bounds().Dx() - wfinal
	}

	// same logic as X, but this time for vertical (Y) axis
	y := faceCenter.Y - (hfinal / 2)
	if y < 0 {
		y = 0
	} else if (y + hfinal) > img.Bounds().Dy() {
		y = img.Bounds().Dy() - hfinal
	}

	// Now crop the image with anchor which is the left-top coordinate of the cropped image
	anchor := image.Pt(x, y)
	config := cutter.Config{
		Width:  wfinal,
		Height: hfinal,
		Anchor: anchor,
		Mode:   cutter.TopLeft,
	}

	croppedImg, err := cutter.Crop(img, config)

	log.Printf("Log config %+v", config)

	// Now resize the cropped square image to target size
	resizedImage := resize.Resize(uint(wfinal), uint(hfinal), croppedImg, resize.Lanczos3)
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
