package format

import (
	"image"
	"image/png"
	"io"
	"log"

	"github.com/disintegration/imageorient"
)

type PngFormat struct {
	id  string
	ext string
}

func (pf *PngFormat) Id() string {
	return pf.id
}

func (pf *PngFormat) Ext() string {
	return pf.ext
}

func (pf *PngFormat) Decode(f io.Reader) (image.Image, error) {
	img, _, err := imageorient.Decode(f)
	if err != nil {
		log.Printf("imageorient.Decode failed: %v\n", err)
		img, err = png.Decode(f)
		if err != nil {
			log.Printf("Error converting from PNG: %V", err)
			return nil, err
		}
	}
	return img, err
}

func (pf *PngFormat) Encode(i image.Image, dest io.Writer, q int) {
	pngenc := png.Encoder{}
	pngenc.CompressionLevel = png.BestCompression
	err := pngenc.Encode(dest, i)
	if err != nil {
		panic("Error converting to PNG")
	}
}
