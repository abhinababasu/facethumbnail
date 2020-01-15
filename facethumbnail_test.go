package facethumbnail

import (
	"image"
	"os"
	"path"
	"testing"
)

func TestDetectFacesInImageFile(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagetall.jpg")
	cascadeFile := path.Join(pwd, "test", "facefinder")

	fd := GetFaceDetector(cascadeFile, -1, -1)
	faces, err := fd.DetectFacesInImageFile(source)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if len(faces) != 1 {
		t.Fatal("Only one face should be found")
	}

	face := faces[0]

	// expected face feature location in test image
	lEye := image.Point{330, 150}
	rEye := image.Point{400, 170}
	nose := image.Point{360, 190}
	mouth := image.Point{360, 220}

	// validate detected face contains these aspects
	if !lEye.In(face) || !rEye.In(face) || !nose.In(face) || !mouth.In(face) {
		t.Errorf("Detected face %v is not correct", face)
	}
}

func TestResizeImageTall(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagetall.jpg")
	cascadeFile := path.Join(pwd, "test", "facefinder")
	destination := testOutputPath("testimagetall_thumb.jpg")

	fd := GetFaceDetector(cascadeFile, -1, -1)

	err := ResizeImage(fd, source, destination, 200)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// TODO: Add validation that thumbnail indeed has the face in it. Run face detection on thumbnail again?
	t.Logf("Check generated file %v", destination)
}

func TestResizeImageWide(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagewide.jpg")
	cascadeFile := path.Join(pwd, "test", "facefinder")
	destination := testOutputPath("testimagewide_thumb.jpg")

	fd := GetFaceDetector(cascadeFile, -1, -1)

	err := ResizeImage(fd, source, destination, 200)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	// TODO: Add validation that thumbnail indeed has the face in it. Run face detection on thumbnail again?
	t.Logf("Check generated file %v", destination)
}

func testOutputPath(fileName string) string {
	pwd, _ := os.Getwd()
	destination := path.Join(pwd, "testoutput")

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		os.MkdirAll(destination, os.ModeDir)
	}

	return path.Join(destination, fileName)
}
