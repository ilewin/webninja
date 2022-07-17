// TO DO: Refactor Cleaner

package services

import (
	"os"
	"time"
)

func Clean(p string) {
	time.AfterFunc(5*time.Minute, func() {
		os.RemoveAll(p)
	})
}
