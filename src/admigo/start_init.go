package main

import (
	"admigo/config"
	"admigo/lang"
	"admigo/models/mcom"
)

func runInit() {
	config.LoadConfig()
	lang.LoadMessages()
	config.SetCsrf()
	mcom.LoadDriver()
}
