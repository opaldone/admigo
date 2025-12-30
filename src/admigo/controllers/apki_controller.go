package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"math"
	"net/http"
	"os"

	"github.com/loovien/ipapk"
	_ "github.com/loovien/ipapk"

	"github.com/julienschmidt/httprouter"
)

const (
	APPS_PT string = "/static/apki"
	MN_APKI string = "apki"
)

func doApkiWebError(w http.ResponseWriter, r *http.Request, err error) {
	WebError(w, r, err, MN_APKI)
}

func humSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	const unit = 1024
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(unit)))

	if i >= len(sizes) {
		i = len(sizes) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/math.Pow(unit, float64(i)), sizes[i])
}

func imgToB64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)

	if err != nil {
		return "", err
	}

	imgBytes := buf.Bytes()
	estring := base64.StdEncoding.EncodeToString(imgBytes)

	return estring, nil
}

func infoApk(pt string, nm string, md string) (ret map[string]any) {
	apa := fmt.Sprintf("%s/%s", pt, nm)
	ai, _ := ipapk.NewAPKParser(apa)

	ret = map[string]any{}

	ret["lbl"] = ai.Label
	ret["vn"] = ai.VersionName
	ret["vc"] = ai.VersionCode
	ret["sz"] = humSize(ai.Size)
	ret["ico"], _ = imgToB64(ai.Icon)
	ret["hre"] = fmt.Sprintf("%s/%s", APPS_PT, nm)
	ret["md"] = md

	return
}

func ApkiIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	wd, err := os.Getwd()
	if err != nil {
		doApkiWebError(w, r, err)
		return
	}

	pt := fmt.Sprintf("%s%s", wd, APPS_PT)

	files, err := os.ReadDir(pt)
	if err != nil {
		doApkiWebError(w, r, err)
		return
	}

	info := []map[string]any{}

	for _, fl := range files {
		fo, err := fl.Info()

		if err != nil {
			doApkiWebError(w, r, err)
			return
		}

		nm := fo.Name()
		ia := infoApk(pt, nm, fo.ModTime().Format("2006-01-02"))
		info = append(info, ia)
	}

	setFrontContent(w, r, "apki", info, nil,
		"apki/ix/index",
	)
}
