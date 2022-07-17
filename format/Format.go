package format

import (
	"image"
	"io"
)

type Format interface {
	Id() string
	Ext() string
	Encode(image.Image, io.Writer, int)
	Decode(io.Reader) (image.Image, error)
}

type Formats struct {
	formats map[string]*Format
}

func (fs *Formats) Register(f Format) {
	fs.formats[f.Id()] = &f
}

func (fs *Formats) Get(id string) *Format {
	if f, ok := fs.formats[id]; ok {
		return f
	}
	panic("Unknown Format")
}

func InitFormats() *Formats {
	fs := &Formats{make(map[string]*Format)}

	jf := &JpegFormat{"image/jpeg", "jpg"}
	pf := &PngFormat{"image/png", "png"}
	gf := &GifFormat{"image/gif", "gif"}
	wf := &WebpFormat{"image/webp", "webp"}

	fs.Register(jf)
	fs.Register(pf)
	fs.Register(gf)
	fs.Register(wf)

	return fs
}
