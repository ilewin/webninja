package format

import (
	"image"
	"image/gif"
	"io"
	"log"

	"github.com/disintegration/imageorient"
)

type GifFormat struct {
	id  string
	ext string
}

func (gf *GifFormat) Id() string {
	return gf.id
}

func (gf *GifFormat) Ext() string {
	return gf.ext
}

func (gf *GifFormat) Decode(f io.Reader) (image.Image, error) {
	img, _, err := imageorient.Decode(f)
	if err != nil {
		log.Printf("imageorient.Decode failed: %v", err)
		img, err = gif.Decode(f)
		if err != nil {
			log.Printf("Error decoding GIF: %v", err)
			return nil, err
		}
	}
	return img, err
}

func (gf *GifFormat) Encode(i image.Image, dest io.Writer, q int) {
	err := gif.Encode(dest, i, nil)
	if err != nil {
		panic("Error converting to JPEG")
	}
}
