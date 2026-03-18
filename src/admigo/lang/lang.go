// Package lang
package lang

import (
	"encoding/json"
	"fmt"
	"os"

	"admigo/applog"
	"admigo/config"
)

const trFolder = "tra"

var lame map[string]string

func readJSON(fileName string) (ret *[]byte) {
	lang := config.Env(false).Lang
	lap := fmt.Sprintf("%s/%s/%s/%s.json", config.Env(false).Static, trFolder, lang, fileName)

	_, err := os.Stat(lap)
	if err != nil {
		return
	}

	laf, err := os.ReadFile(lap)
	if err != nil {
		applog.Danger("Cannot open file", err)
		return
	}

	ret = &laf

	return
}

func LoadMessages() {
	pcont := readJSON("messages")

	if pcont == nil {
		return
	}

	json.Unmarshal(*pcont, &lame)
}

func LoadLabels(fileName string) (ret *map[string]string) {
	pcont := readJSON(fileName)

	if pcont == nil {
		return
	}

	json.Unmarshal(*pcont, &ret)

	return
}

func NeedTra() bool {
	return len(lame) != 0
}

func Re(str ...string) string {
	if len(str) == 0 {
		return ""
	}

	key := str[0]

	if len(lame) == 0 {
		return key
	}

	ret, ok := lame[key]
	if !ok {
		return key
	}

	if len(str) == 2 {
		ret = fmt.Sprintf(ret, str[1])
	}

	return ret
}
