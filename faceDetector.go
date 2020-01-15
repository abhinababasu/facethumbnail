package facethumbnail

import (
	"image"
	"io/ioutil"
	"log"

	pigo "github.com/esimov/pigo/core"
)

// faceDetector struct contains Pigo face detector general settings.
type FaceDetector struct {
	angle        float64
	cascadeFile  string
	minSize      int
	maxSize      int
	shiftFactor  float64
	scaleFactor  float64
	iouThreshold float64
}

// GetFaceDetector returns an instance of facedetector initialized with a cascadefile with expected minimum and
// maximum size of the face. min and max can be -1 to use default values
func GetFaceDetector(cascadeFile string, minSize, maxSize int) *FaceDetector {

	if minSize <= 0 {
		minSize = 20
	}

	if maxSize <= 0 {
		maxSize = 1000
	}

	fd := &FaceDetector{
		angle:        0.0,
		cascadeFile:  cascadeFile,
		minSize:      minSize,
		maxSize:      maxSize,
		shiftFactor:  0.1,
		scaleFactor:  1.1,
		iouThreshold: 0.2,
	}

	return fd
}

// DetectFacesInImageFile detect faces in a image file
func (fd *FaceDetector) DetectFacesInImageFile(sourceFilePath string) ([]image.Rectangle, error) {

	faces, err := fd.detectFaces(sourceFilePath)
	if err != nil {
		log.Fatalf("Detection error: %v", err)
	}

	rects := fd.generateFaceRects(faces)

	if err != nil {
		log.Fatalf("Error creating the image output: %s", err)
	}

	return rects, nil
}

// detectFaces run the detection algorithm over the provided source image.
func (fd *FaceDetector) detectFaces(source string) ([]pigo.Detection, error) {
	src, err := pigo.GetImage(source)
	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	imgParams := &pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}

	cParams := pigo.CascadeParams{
		MinSize:     fd.minSize,
		MaxSize:     fd.maxSize,
		ShiftFactor: fd.shiftFactor,
		ScaleFactor: fd.scaleFactor,
		ImageParams: *imgParams,
	}

	cascadeFile, err := ioutil.ReadFile(fd.cascadeFile)
	if err != nil {
		return nil, err
	}

	p := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := p.Unpack(cascadeFile)
	if err != nil {
		return nil, err
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := classifier.RunCascade(cParams, fd.angle)

	// Calculate the intersection over union (IoU) of two clusters.
	faces = classifier.ClusterDetections(faces, fd.iouThreshold)

	return faces, nil
}

// drawFaces marks the detected faces with a rectangle.
func (fd *FaceDetector) generateFaceRects(faces []pigo.Detection) []image.Rectangle {
	var qThresh float32 = 5.0
	var rects []image.Rectangle

	for _, face := range faces {
		if face.Q > qThresh {
			rect := image.Rect(
				face.Col-face.Scale/2,
				face.Row-face.Scale/2,
				face.Col+face.Scale/2,
				face.Row+face.Scale/2,
			)

			rects = append(rects, rect)
		}
	}

	return rects
}
