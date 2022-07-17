package abstract

import (
	"image"
	"io"
)

type Decoder interface {
	Decode(io.Reader) image.Image
}
