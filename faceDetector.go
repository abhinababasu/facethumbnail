package facethumbnail

import (
	"fmt"
	"image"
	"io/ioutil"

	pigo "github.com/esimov/pigo/core"
)

// FDetect struct contains Pigo face detector general settings.
type FDetect struct {
	angle        float64
	cascadeFile  string
	minSize      int
	maxSize      int
	shiftFactor  float64
	scaleFactor  float64
	iouThreshold float64
	classifier   *pigo.Pigo
	initialized bool
}

// GetFaceDetector returns an instance of facedetector with path to binary cascade file
func GetFaceDetector(cascadeFile string) *FDetect {
	fd := &FDetect{}
	fd.cascadeFile = cascadeFile
	return fd
}

// Init initializes facedetector with a cascadefile with expected minimum and maximum size of the face.
// min and max can be -1 to use default values
func (fd *FDetect) Init(minSize, maxSize int) error {
	if fd.initialized {
		return fmt.Errorf("Already initialized")
	}

	if minSize <= 0 {
		minSize = 20
	}

	if maxSize <= 0 {
		maxSize = 1000
	}

	fd.angle = 0.0

	fd.minSize = minSize
	fd.maxSize = maxSize
	fd.shiftFactor = 0.1
	fd.scaleFactor = 1.1
	fd.iouThreshold = 0.2

	cascadeFile, err := ioutil.ReadFile(fd.cascadeFile)
	if err != nil {
		return err
	}

	p := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	fd.classifier, err = p.Unpack(cascadeFile)
	if err != nil {
		return err
	}

	fd.initialized = true

	return nil
}

// DetectFacesInImageFile detect faces in a image file
func (fd *FDetect) DetectFacesInImageFile(sourceFilePath string) ([]image.Rectangle, error) {

	faces, err := fd.detectFaces(sourceFilePath)
	if err != nil {
		return nil, err
	}

	rects := fd.generateFaceRects(faces)

	return rects, nil
}

// detectFaces run the detection algorithm over the provided source image.
func (fd *FDetect) detectFaces(source string) ([]pigo.Detection, error) {
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

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := fd.classifier.RunCascade(cParams, fd.angle)

	// Calculate the intersection over union (IoU) of two clusters.
	faces = fd.classifier.ClusterDetections(faces, fd.iouThreshold)

	return faces, nil
}

// drawFaces marks the detected faces with a rectangle.
func (fd *FDetect) generateFaceRects(faces []pigo.Detection) []image.Rectangle {
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
