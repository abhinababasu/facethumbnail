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

	fd := GetFaceDetector(cascadeFile)
	fd.Init(-1, -1)
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
	destination := testOutputPath("testimagetall_thumb.jpg")

	runTestImage(source, destination, 1, t)
}

func TestResizeImageWide(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagewide.jpg")
	destination := testOutputPath("testimagewide_thumb.jpg")

	runTestImage(source, destination, 1, t)
}

func TestResizeImageManyPeople(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagemany.jpg")
	destination := testOutputPath("testimagemany_thumb.jpg")

	runTestImage(source, destination, 12, t)
}


func TestResizeImageNoPeople(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagenone.jpg")
	destination := testOutputPath("testimagenone_thumb.jpg")

	runTestImage(source, destination, 0, t)
}

func TestResizeImageNoFaceDetection(t *testing.T) {
	pwd, _ := os.Getwd()
	source := path.Join(pwd, "test", "testimagetall.jpg")
	destination := testOutputPath("testimagemnodetection_thumb.jpg")
	result, _ := ResizeImage(nil, source, destination, 200)

	if result.FacesCount != 0 {
		t.Errorf("Expected face count 0 did not match actual %v", result.FacesCount)
	}
}

func testOutputPath(fileName string) string {
	pwd, _ := os.Getwd()
	destination := path.Join(pwd, "testoutput")

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		os.MkdirAll(destination, os.ModeDir)
	}

	return path.Join(destination, fileName)
}

func runTestImage(srcPath, dstPath string, expectedFaceCount int, t *testing.T) {
	pwd, _ := os.Getwd()
	cascadeFile := path.Join(pwd, "test", "facefinder")

	fd := GetFaceDetector(cascadeFile)
	fd.Init(-1, -1)

	result, err := ResizeImage(fd, srcPath, dstPath, 200)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if result.FacesCount != expectedFaceCount {
		t.Errorf("Expected face count %v did not match actual %v", expectedFaceCount, result.FacesCount)
	}
	// TODO: Add validation that thumbnail indeed has the face in it. 
	// Run face detection on thumbnail again?
	t.Logf("Check generated file %v", dstPath)

}