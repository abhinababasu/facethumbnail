# facethumbnail
Thumbnail generator that preserves the face in the final image

Consider the image below

![Sample Image](./samples/image.jpg "Sample Image")

Now if we simply generate a square thumbnail from this, it will generate something like below

![Sample Thumbnail](samples/thumbnail_bad.jpg "Sample Thumbnail")

The head is cut off. **facethumbnai**l tries to detect faces in pictures and attempts to choose a region of 
image so that the thumbnail contains the face properly, as follows

![Sample Thumbnail](samples/thumbnail_good.jpg "Sample Thumbnail")


## Build/Test
Clone this repo and in the folder run the following commands

```
go get
go build .
go test
```

