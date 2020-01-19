# facethumbnail
Thumbnail generator that preserves the face in the final image

Consider the image below

![Sample Image](./samples/image.jpg "Sample Image")

Now if we simply generate a square thumbnail from this, it will generate something like below

![Sample Thumbnail](samples/thumbnail_bad.jpg "Sample Thumbnail")

The head is cut off. **facethumbnai**l tries to detect faces in pictures and attempts to choose a region of 
image so that the thumbnail contains the face properly, as follows

![Sample Thumbnail](samples/thumbnail_good.jpg "Sample Thumbnail")

Facethumbnail also supports images with multiple faces and uses the largest face as the center of the thumbnail.

## Build/Test
Clone this repo and in the folder run the following commands

```
go get
go build .
go test . -v
```

## Sample
A sample use is
``` go
pwd, _ := os.Getwd()                                   // get the current directory
cascadeFile := path.Join(pwd, "test", "facefinder")    // get path to the facefinder binary
fd := GetFaceDetector(cascadeFile, -1, -1)             // create facedetector instance
result, err := ResizeImage(fd, srcPath, dstPath, 200)  // resize image!!
```

For more sample usage see the [test file](facethumbnail_test.go)

## How Does this work
This project uses github.com/esimov/pigo/core to detect the face, then generate thumbnail with the detected face centered.

## Known issues
1. Works for only jpeg