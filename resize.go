package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

func resizeImage(img io.Reader, maxX int, maxY int) (io.Reader, error) {

	i, t, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	i = resize.Thumbnail(uint(maxX), uint(maxY), i, resize.Lanczos3)

	r := bytes.NewBuffer(nil)

	switch t {
	case "jpeg":

		err = jpeg.Encode(r, i, nil)
		if err != nil {
			return nil, err
		}
	case "png":
		err = png.Encode(r, i)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
