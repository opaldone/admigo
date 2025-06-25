package mcom

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"admigo/common"
	"admigo/config"
)

const MAX_SIZE = 3 * 1024 * 1024

func SaveSomeFile(r *http.Request, name string, path string) (fileName string, err error) {
	file, info, err := r.FormFile(name)
	if err != nil {
		return "", nil
	}

	defer file.Close()

	size := info.Size

	if size > MAX_SIZE {
		err = fmt.Errorf("the file size is limited to %d", MAX_SIZE)
		return
	}

	fb, err := io.ReadAll(file)
	if err != nil {
		return
	}

	ftype := http.DetectContentType(fb)
	ext := ""
	switch ftype {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpeg"
	case "image/jpg":
		ext = ".jpg"
	case "image/gif":
		ext = ".gif"
	default:
		err = fmt.Errorf("not allowed extension %s", ftype)
		return
	}

	rfn := common.RandUID()

	fn := rfn + ext
	pt := fmt.Sprintf("./%s/images/%s/%s", config.Env(false).Static, path, fn)

	cf, err := os.Create(pt)
	if err != nil {
		return
	}

	_, err = cf.Write(fb)
	if err != nil {
		return
	}

	err = cf.Close()
	if err != nil {
		return
	}

	fileName = fn

	return
}
