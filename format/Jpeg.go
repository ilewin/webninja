package format

import (
	"image"
	"image/jpeg"
	"io"
	"log"

	"github.com/disintegration/imageorient"
)

type JpegFormat struct {
	id  string
	ext string
}

func (jf *JpegFormat) Id() string {
	return jf.id
}

func (jf *JpegFormat) Ext() string {
	return jf.ext
}

func (jf *JpegFormat) Decode(f io.Reader) (image.Image, error) {
	img, _, err := imageorient.Decode(f)
	if err != nil {
		log.Printf("imageorient.Decode failed: %v", err)
		img, err = jpeg.Decode(f)
		if err != nil {
			log.Printf("Error converting from JPEG: %v", err)
			return nil, err
		}
	}

	return img, err

}

func (jf *JpegFormat) Encode(i image.Image, dest io.Writer, q int) {
	var opt jpeg.Options
	opt.Quality = q
	err := jpeg.Encode(dest, i, &opt)
	if err != nil {
		log.Printf("Error converting to JPEG")
	}
}
