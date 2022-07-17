package abstract

import (
	"image"
	"io"
)

type Encoder interface {
	Encode(image.Image, io.Writer, int)
}
