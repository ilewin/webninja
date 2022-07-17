package format

import (
	"image"
	"io"
	"log"

	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type WebpFormat struct {
	id  string
	ext string
}

func (wf *WebpFormat) Id() string {
	return wf.id
}

func (wf *WebpFormat) Ext() string {
	return wf.ext
}

func (wf *WebpFormat) Decode(f io.Reader) (image.Image, error) {
	img, err := webp.Decode(f, &decoder.Options{})
	if err != nil {
		log.Printf("Error converting from GIF: %v", err)
		return nil, err
	}
	return img, err
}

func (wf *WebpFormat) Encode(i image.Image, dest io.Writer, q int) {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(q)) //NewLosslessEncoderOptions(encoder.PresetPhoto, 1)
	if err != nil {
		log.Fatalln(err)
	}

	if err := webp.Encode(dest, i, options); err != nil {
		log.Fatalln(err)
	}

}
