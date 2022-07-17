package handlers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"webp.ninja/services"
	"webp.ninja/utils"
)

type WebFile struct {
	// GetName() string
	// GetSize() int64
	// GetFormat() string
	// EncodeTo() string
	// GetPath() string
	name     string
	size     int64
	format   string
	encodeTo string
	path     string
}

func (wf WebFile) GetName() string {
	return wf.name
}

func (wf WebFile) GetSize() int64 {
	return wf.size
}

func (wf WebFile) GetFormat() string {
	return wf.format
}

func (wf WebFile) EncodeTo() string {
	return wf.encodeTo
}

func (wf WebFile) GetPath() string {
	return wf.path
}

func ConvertHandler(c *fiber.Ctx) error {

	config := utils.GetConfig()

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File[config.App_Files_Field]

	sid := uuid.New().String()

	os.Mkdir(fmt.Sprintf(config.App_Storage+"%s", sid), 0755)
	defer services.Clean(fmt.Sprintf(config.App_Storage+"%s", sid))

	var wfiles []services.InFile

	for _, f := range files {

		lp := fmt.Sprintf(config.App_Storage+"%s/%s", sid, f.Filename)
		err := c.SaveFile(f, lp)
		if err != nil {
			return err
		}

		frm, err := utils.GetFormat(form.Value["convertTo"][0])

		if err != nil {
			continue
		}

		wfiles = append(wfiles, WebFile{
			name:     f.Filename,
			size:     f.Size,
			format:   f.Header.Values("Content-Type")[0],
			encodeTo: frm,
			path:     lp,
		})

	}

	p := services.HandleConvert(wfiles, config, sid)

	go services.UpdateMeta(&p.Files)

	return c.JSON(p.Files)
}
