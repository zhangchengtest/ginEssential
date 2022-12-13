package util

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"testing"
)

func getTestImage() (image.Image, error) {
	f, err := os.Open("d:/test/sss.jpg")
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func TestSplitImageWithIterator(t *testing.T) {
	in, err := getTestImage()
	if err != nil {
		t.Fatalf("failed to load test image: %s", err)
	}

	// set cfg
	cfg := Config{X: 3, Y: 3}

	// split
	it, err := SplitImageWithIterator(in, cfg)
	if err != nil {
		t.Fatalf("unexpected failure: %s", err)
	}

	// drain the parts to check expected count
	outs := ConsumeIterator(it)

	for i, tt := range outs {
		out, _ := os.Create("d:/test/sub/img" + strconv.Itoa(i) + ".png")
		defer out.Close()

		err = png.Encode(out, tt)
		//jpeg.Encode(out, img, nil)
		if err != nil {
			log.Println(err)
		}
	}

	if len(outs) != 9 {
		t.Fatalf("wrong number of parts: %d", len(outs))
	}
}
