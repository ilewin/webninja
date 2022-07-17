package utils

import (
	"errors"
	"strings"
)

func NewFileName(fname string, ext string) string {
	fnpcs := strings.Split(fname, ".")
	fnpcs[len(fnpcs)-1] = ext
	return strings.Join(fnpcs, ".")
}

func GetFormat(ff string) (string, error) {
	fmts := map[string]string{
		"WEBP":       "image/webp",
		"JPG":        "image/jpeg",
		"PNG":        "image/png",
		"GIF":        "image/gif",
		"image/webp": "image/webp",
		"image/jpeg": "image/jpeg",
		"image/png":  "image/png",
		"image/gif":  "image/gif",
	}
	if inf, ok := fmts[ff]; ok {
		return inf, nil
	}

	return "", errors.New("Can't find requsted image format")
}
