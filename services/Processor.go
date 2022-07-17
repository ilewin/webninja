package services

import (
	"fmt"
	"log"
	"os"
	"sync"

	"webp.ninja/format"
	"webp.ninja/utils"
)

var (
	wg   = sync.WaitGroup{}
	lock = sync.Mutex{}
)

type File struct {
	Sid      string `json:"sid"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	NSize    int64  `json:"newsize"`
	Format   string `json:"format"`
	EncodeTo string `json:"encodeto"`
	Error    string `json: "error"`
	Path     string
}

func (f *File) Convert(from format.Format, to format.Format) { // REFACTOR

	config := utils.GetConfig()

	file, err := os.Open(f.Path)
	if err != nil {
		log.Fatalln(err)
	}

	iimg, err := from.Decode(file)

	if err != nil {
		f.Error = err.Error()
		wg.Done()
		return
	}

	nfp := utils.NewFileName(f.Path, to.Ext())

	oimg, err := os.Create(nfp)
	if err != nil {
		log.Fatal(err)
	}
	defer oimg.Close()

	to.Encode(iimg, oimg, config.Compression)

	f.Name = utils.NewFileName(f.Name, to.Ext())
	f.Path = nfp
	f.Format = to.Id()
	if stat, err := oimg.Stat(); err == nil {
		f.NSize = stat.Size()
	}

	wg.Done()
}

type Processor struct {
	Files   []File
	formats *format.Formats
}

func (p *Processor) Convert() {
	wg.Add(len(p.Files))
	for i := 0; i < len(p.Files); i++ {
		go p.Files[i].Convert(*p.formats.Get(p.Files[i].Format), *p.formats.Get(p.Files[i].EncodeTo))
	}
	wg.Wait()
}

func NewProcessor(l int) *Processor {
	fs := format.InitFormats()
	return &Processor{make([]File, 0, l), fs}
}

type InFile interface {
	GetName() string
	GetSize() int64
	GetFormat() string
	EncodeTo() string
	GetPath() string
}

func HandleConvert(files []InFile, config *utils.Config, sid string) *Processor {

	if len(files) == 0 {
		return nil
	}
	p := NewProcessor(len(files))

	for _, f := range files {

		frm, err := utils.GetFormat(f.EncodeTo())
		fmt.Println(f.GetFormat(), f.EncodeTo())
		if err != nil {
			continue
		}

		p.Files = append(p.Files, File{
			Sid:      sid,
			Name:     f.GetName(),
			Size:     f.GetSize(),
			NSize:    0,
			Format:   f.GetFormat(),
			EncodeTo: frm,
			Path:     f.GetPath(),
		})

	}

	p.Convert()

	return p

}
